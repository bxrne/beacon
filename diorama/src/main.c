#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "freertos/timers.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "button_task.h"
#include "traffic_light_task.h"
#include "cJSON.h"

QueueHandle_t event_queue;

void app_main(void)
{
    esp_rom_gpio_pad_select_gpio(CAR_GREEN_PIN);
    esp_rom_gpio_pad_select_gpio(CAR_YELLOW_PIN);
    esp_rom_gpio_pad_select_gpio(CAR_RED_PIN);
    esp_rom_gpio_pad_select_gpio(PED_GREEN_PIN);
    esp_rom_gpio_pad_select_gpio(PED_RED_PIN);
    esp_rom_gpio_pad_select_gpio(PED_BUTTON_PIN);

    gpio_set_direction(CAR_GREEN_PIN, GPIO_MODE_OUTPUT);
    gpio_set_direction(CAR_YELLOW_PIN, GPIO_MODE_OUTPUT);
    gpio_set_direction(CAR_RED_PIN, GPIO_MODE_OUTPUT);
    gpio_set_direction(PED_GREEN_PIN, GPIO_MODE_OUTPUT);
    gpio_set_direction(PED_RED_PIN, GPIO_MODE_OUTPUT);
    gpio_set_direction(PED_BUTTON_PIN, GPIO_MODE_INPUT);

    // Enable internal pull-up resistor for the pedestrian button pin
    gpio_set_pull_mode(PED_BUTTON_PIN, GPIO_PULLUP_ONLY);
    ESP_LOGI(TAG, "[DEBUG] Button pin configured with pull-up resistor");

    event_queue = xQueueCreate(10, sizeof(event_t));

    if (event_queue == NULL)
    {
        ESP_LOGE(TAG, "[ERROR] Failed to create event queue");
        return;
    }

    if (xTaskCreate(traffic_light_task, "Traffic Light Task", 2048, NULL, 1, NULL) != pdPASS)
    {
        ESP_LOGE(TAG, "[ERROR] Failed to create Traffic Light Task");
    }

    if (xTaskCreate(button_task, "Button Task", 2048, NULL, 1, NULL) != pdPASS)
    {
        ESP_LOGE(TAG, "[ERROR] Failed to create Button Task");
    }

    // Initial traffic light state
    gpio_set_level(CAR_GREEN_PIN, 1);
    gpio_set_level(CAR_YELLOW_PIN, 0);
    gpio_set_level(CAR_RED_PIN, 0);
    gpio_set_level(PED_GREEN_PIN, 0);
    gpio_set_level(PED_RED_PIN, 1);
    ESP_LOGI(TAG, "[DEBUG] Initial traffic light state set");
}
