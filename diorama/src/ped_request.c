#include <freertos/FreeRTOS.h>
#include "ped_request.h"
#include "driver/gpio.h"
#include "freertos/queue.h"
#include "config.h"
#include "globals.h"
#include "freertos/task.h"
#include "esp_log.h"

#define TAG "PED_REQUEST"

void button_isr_handler(void *arg)
{
  uint32_t button_pressed = 1;
  xQueueSendFromISR(pedestrianRequestQueue, &button_pressed, NULL);
}

void init_ped_request(void)
{
  // Initialize pedestrian button GPIO
  esp_rom_gpio_pad_select_gpio(PED_BUTTON_PIN);
  gpio_set_direction(PED_BUTTON_PIN, GPIO_MODE_INPUT);
  gpio_set_intr_type(PED_BUTTON_PIN, GPIO_INTR_POSEDGE);

  gpio_install_isr_service(0);
  gpio_isr_handler_add(PED_BUTTON_PIN, button_isr_handler, NULL);
}

void PedestrianRequestTask(void *pvParameters)
{
  uint32_t button_pressed;
  while (true)
  {
    if (xQueueReceive(pedestrianRequestQueue, &button_pressed, portMAX_DELAY))
    {
      ESP_LOGI(TAG, "Pedestrian button pressed");
      xSemaphoreGive(xPedestrianSemaphore); // Give traffic light task permission to change the light
    }
  }
}
