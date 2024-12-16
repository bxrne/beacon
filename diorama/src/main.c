#include <freertos/FreeRTOS.h>
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>

#include "freertos/task.h"
#include "freertos/queue.h"
#include "freertos/semphr.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "car_light.h"
#include "ped_light.h"
#include "ped_request.h"
#include "globals.h"

#define TAG "MAIN"

QueueHandle_t pedestrianRequestQueue;
SemaphoreHandle_t xPedestrianSemaphore;

void init_gpio(void)
{
    ESP_LOGI(TAG, "Initializing GPIOs");

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

    gpio_set_intr_type(PED_BUTTON_PIN, GPIO_INTR_POSEDGE);
    gpio_install_isr_service(0);
    gpio_isr_handler_add(PED_BUTTON_PIN, button_isr_handler, NULL);

    gpio_set_level(PED_RED_PIN, 1); // Ensure pedestrian light is red by default
}

void app_main(void)
{
    ESP_LOGI(TAG, "Starting application");

    pedestrianRequestQueue = xQueueCreate(10, sizeof(uint32_t));
    xPedestrianSemaphore = xSemaphoreCreateBinary();

    init_gpio();
    init_ped_request();

    xTaskCreate(CarLightTask, "CarLightTask", 2048, NULL, 1, NULL);
    xTaskCreate(PedestrianLightTask, "PedestrianLightTask", 2048, NULL, 1, NULL);
}
