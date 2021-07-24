package main

import (
	"context"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/shauncampbell/zigbee2mqtt-mongodb-broker/internal/config"
	"github.com/shauncampbell/zigbee2mqtt-mongodb-broker/pkg/zigbee2mqtt"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var rootCmd = &cobra.Command{
	Use:  "zigbee2mqtt-mongodb-broker",
	RunE: runBroker,
}

const defaultTimeout = 10 * time.Second

func runBroker(cmd *cobra.Command, args []string) error {
	cfg, err := config.Read()

	if err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	if cfg.MongoDB.Host == "" {
		return fmt.Errorf("the environment variable BROKER_MONGODB_HOST must be specified")
	}

	// connect to mongodb
	var mongoDBURI string
	if cfg.MongoDB.URI == "" {
		if cfg.MongoDB.Username == "" {
			mongoDBURI = fmt.Sprintf("mongodb://%s:%d/%s", cfg.MongoDB.Host, cfg.MongoDB.Port, cfg.MongoDB.Database)
		} else {
			mongoDBURI = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
				cfg.MongoDB.Username,
				cfg.MongoDB.Password,
				cfg.MongoDB.Host,
				cfg.MongoDB.Port,
				cfg.MongoDB.Database)
		}
	} else {
		mongoDBURI = cfg.MongoDB.URI
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	mgoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		return fmt.Errorf("failed to connect to mongodb: %w", err)
	}
	defer func() {
		if err = mgoClient.Disconnect(ctx); err != nil {
			log.Error().Msgf("failed to disconnect from mongodb: %s", err.Error())
		}
	}()
	// Ping the primary
	if err := mgoClient.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to connect to mongodb primary: %w", err)
	}

	handler := zigbee2mqtt.New(mgoClient)

	mqttOptions := mqtt.NewClientOptions()
	mqttOptions.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.MQTT.Host, cfg.MQTT.Port))
	if cfg.MQTT.Username != "" {
		mqttOptions.Username = cfg.MQTT.Username
	}
	if cfg.MQTT.Password != "" {
		mqttOptions.Password = cfg.MQTT.Password
	}
	mqttOptions.SetClientID("zigbee2mqtt-mongodb-broker")
	mqttOptions.SetDefaultPublishHandler(handler.DefaultPublished)
	mqttOptions.OnConnect = handler.Connected
	mqttOptions.OnConnectionLost = handler.Disconnected
	client := mqtt.NewClient(mqttOptions)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for {
		time.Sleep(defaultTimeout)
		if !client.IsConnected() {
			log.Error().Msg("connection to mqtt was severed")
			break
		}
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error().Msgf("unable to run broker: %s", err.Error())
		if _, e := fmt.Fprintf(os.Stderr, "unable to run broker: %s\n", err.Error()); e != nil {
			log.Error().Msgf("unable to write to stderr: %s", e.Error())
			log.Error().Msgf("unable to run broker: %s", err.Error())
		}
		os.Exit(1)
	}
	os.Exit(0)
}
