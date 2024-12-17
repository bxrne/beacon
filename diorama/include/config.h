#ifndef CONFIG_H
#define CONFIG_H

// Wi-Fi configuration
#define WIFI_SSID "coldspot"
#define WIFI_PASS "helloworld1"

// Traffic light GPIO pins and durations
#define CAR_GREEN_PIN 18
#define CAR_YELLOW_PIN 19
#define CAR_RED_PIN 21

#define CAR_GREEN_DURATION 5000 // Milliseconds
#define CAR_YELLOW_DURATION 1000
#define CAR_RED_DURATION 3000

// Pedestrian light GPIO pins and durations
#define PED_GREEN_PIN 22
#define PED_RED_PIN 23
#define PED_GREEN_DURATION 7000
#define PED_BUTTON_PIN 25 // Button to request pedestrian light

#endif // CONFIG_H
