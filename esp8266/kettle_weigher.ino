#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <WiFiClient.h>
#include <HX711.h>
#include <my_secrets.h> // Library with SECRET_SSID, SECRET_WIFI_PASSWORD

const uint8_t DOUT_PIN = 4;               // DOUT pin for HX711
const uint8_t SCK_PIN = 5;                // SCK pin for HX711
const float CALIBRATION_FACTOR = -24.83;  // Scale calibration factor
const float MAX_WATER_ML = 1700;          // The maximum amount of water the kettle can hold in ml
const uint8_t MEASUREMENT_AMOUNT = 10;    // The number of measurements to take before calculating average
const uint8_t WATER_LEVEL_ERROR = 5;      // The allowable error of the water_level that doesn't require sending an update

// Inserting Wifi creds
const char* SSID = SECRET_SSID;
const char* WIFI_PASSWORD = SECRET_WIFI_PASSWORD;

const char* KETTLE_WEIGHER_URL = "https://kettle-weigher.alexdenisova.ru/v1.0/user/device/state";

const uint8_t DELAY_SEC = 60;  // Time in between measurements

HX711 scale;
float previous_water_level;

void setup() {
  Serial.begin(115200);

  scale.begin(DOUT_PIN, SCK_PIN);
  scale.set_scale(CALIBRATION_FACTOR);
  scale.tare();
  Serial.println("Scale tared");

  WiFi.begin(SSID, WIFI_PASSWORD);
  Serial.print("Connecting");
  while (WiFi.status() != WL_CONNECTED) {
    Serial.print(".");
    delay(500);
  }
  Serial.printf("\nConnected to WiFi network with IP Address: ");
  Serial.println(WiFi.localIP());
}

void loop() {
  if (WiFi.status() == WL_CONNECTED) {
    if (scale.is_ready()) {
      float weight = scale.get_units(MEASUREMENT_AMOUNT);
      uint8_t water_level;
      if (weight <= 0) {
        water_level = 0;
      } else if (weight >= MAX_WATER_ML) {
        water_level = 100;
      } else {
        water_level = weight / MAX_WATER_ML * 100;
      }
      Serial.printf("Measured: %d%% (%.2fg)\n", water_level, weight);
      if ((water_level < previous_water_level - WATER_LEVEL_ERROR) || (previous_water_level + WATER_LEVEL_ERROR < water_level)) {
        std::unique_ptr<BearSSL::WiFiClientSecure> client(new BearSSL::WiFiClientSecure);
        client->setInsecure(); // Ignore SSL certificate validation

        HTTPClient https;
        https.begin(*client, KETTLE_WEIGHER_URL);
        https.addHeader("Content-Type", "application/json");
        String httpRequestData = String("{\"device_id\": \"kettle-weigher\",\"type\": \"property\",\"instance\": \"water_level\",\"value\": ") + water_level + "}";
        Serial.printf("Request body: %s\n", httpRequestData.c_str());

        // Send HTTP PATCH request
        int httpResponseCode = https.PATCH(httpRequestData);
        Serial.printf("HTTP Response code: %d\n", httpResponseCode);

        https.end();
        previous_water_level = water_level;
      }
    } else {
      Serial.println("Scale Disconnected");
    }
  } else {
    Serial.println("WiFi Disconnected");
  }
  delay(DELAY_SEC * 1000);
}