package zigbee2mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Handler struct {
	database *mongo.Database
	logger zerolog.Logger
	devices map[string]*Device
}

func (h *Handler) Connected(client mqtt.Client) {
	token := client.Subscribe("zigbee2mqtt/bridge/devices", 1, h.deviceListChangedHandler)
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

func (h *Handler) DefaultPublished(client mqtt.Client, message mqtt.Message) {

}

func (h *Handler) deviceListChangedHandler(client mqtt.Client, message mqtt.Message) {
	var devices []Device
	err := json.Unmarshal(message.Payload(), &devices)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for _, device := range devices {
			if h.devices[device.IEEEAddress] == nil {
				h.subscribeToDevice(client, &device)
			}
		}
	}
}

func (h *Handler) subscribeToDevice(client mqtt.Client, device *Device) {
	logger := h.logger.With().Str("friendly_name", device.FriendlyName).Str("device_id", device.IEEEAddress).Logger()
	logger.Info().Msgf("subscribing to device")
	client.Subscribe("zigbee2mqtt/"+device.FriendlyName, 1, h.deviceEventHandler(device, logger))
	h.persistDeviceToMongodb(*device)
}

func (h *Handler) deviceEventHandler(device *Device, logger zerolog.Logger) mqtt.MessageHandler {
	d := *device
	return func(client mqtt.Client, message mqtt.Message) {
		logger.Info().Msgf("received new status message")
		var m map[string]interface{}
		err := json.Unmarshal(message.Payload(), &m)
		if err != nil {
			logger.Error().Msgf("failed to unmarshall payload: %s", err.Error())
			return
		}

		h.persistStateToMongoDB(d, m)
	}
}

func (h *Handler) persistDeviceToMongodb(device Device) {
	// Set up mongodb variables
	collection := h.database.Collection("current_state")

	set := bson.D{
		{"friendly_name", device.FriendlyName},
		{ "type", device.Type},
		{ "network_address", device.NetworkAddress},
		{"model", device.Definition.Model},
		{ "vendor", device.Definition.Vendor},
	}

	filter := bson.D{ { "_id", "mqtt_"+device.IEEEAddress }}
	_, err := collection.UpdateOne(context.Background(), filter, bson.D{{"$set", set}}, options.Update().SetUpsert(true))
	if err != nil {
		h.logger.Error().Msgf("failed to persist to mongodb: %s", err.Error())
	}
}

func (h *Handler) persistStateToMongoDB(device Device, state map[string]interface{}) {
	// Figure out the current time
	t := time.Now().UTC()
	y, m, d := time.Now().UTC().Date()
	ts := t.Unix()

	// Set up mongodb variables
	collection := h.database.Collection("current_state")
	set := bson.D{{"update_ts", ts},{"friendly_name", device.FriendlyName}}
	push := bson.D{}

	// Iterate through the changes.
	for _, attr := range device.Definition.Exposes {
		if state[attr.Property] != nil {
			set = append(set, bson.E{attr.Property, state[attr.Property]})
			push = append(push, bson.E{ attr.Property, bson.D{{"ts", ts}, {"v", state[attr.Property]}}})
		}
	}

	filter := bson.D{ { "_id", "mqtt_"+device.IEEEAddress }}
	_, err := collection.UpdateOne(context.Background(), filter, bson.D{{"$set", set}}, options.Update().SetUpsert(true))
	if err != nil {
		h.logger.Error().Msgf("failed to persist to mongodb: %s", err.Error())
	}

	collection = h.database.Collection("historical_state")

	filter = bson.D{
		{ "ieee_address", device.IEEEAddress},
		{ "friendly_name", device.FriendlyName},
		{ "y", y },
		{ "m", m },
		{ "d", d},
	}
	_, err = collection.UpdateOne(context.Background(), filter, bson.D{{"$push", push}}, options.Update().SetUpsert(true))
	if err != nil {
		h.logger.Error().Msgf("failed to persist to mongodb: %s", err.Error())
	}
}


func (h *Handler) Disconnected(client mqtt.Client, err error) {
	fmt.Println("connection failed")
}

func New(client *mongo.Client) *Handler {
	return &Handler{devices: make(map[string]*Device), logger: log.Logger, database: client.Database("home")}
}