#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>

#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/semphr.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "button_task.h"

extern QueueHandle_t event_queue;
extern SemaphoreHandle_t button_semaphore;

static void IRAM_ATTR button_isr_handler(void *arg)
{
  BaseType_t xHigherPriorityTaskWoken = pdFALSE;
  xSemaphoreGiveFromISR(button_semaphore, &xHigherPriorityTaskWoken);
  portYIELD_FROM_ISR(xHigherPriorityTaskWoken);
}

void button_task(void *pvParameters)
{
  ESP_LOGI("BUTTON_TASK", "Button task started");

  // Configure button GPIO with interrupt
  gpio_config_t io_conf = {
      .pin_bit_mask = (1ULL << PED_BUTTON_PIN),
      .mode = GPIO_MODE_INPUT,
      .pull_up_en = GPIO_PULLUP_ENABLE,
      .pull_down_en = GPIO_PULLDOWN_DISABLE,
      .intr_type = GPIO_INTR_NEGEDGE};

  gpio_config(&io_conf);
  gpio_install_isr_service(ESP_INTR_FLAG_DEFAULT);
  gpio_isr_handler_add(PED_BUTTON_PIN, button_isr_handler, NULL);

  while (1)
  {
    // Wait for semaphore from ISR
    if (xSemaphoreTake(button_semaphore, portMAX_DELAY) == pdTRUE)
    {
      // Debounce handling
      vTaskDelay(pdMS_TO_TICKS(DEBOUNCE_TIME_MS));
      if (gpio_get_level(PED_BUTTON_PIN) == 0)
      {
        // Button is still pressed
        event_t event = EVENT_BUTTON_PRESS;
        xQueueSend(event_queue, &event, portMAX_DELAY);
        ESP_LOGI("BUTTON_TASK", "Button press event sent");
      }
    }
  }
}
