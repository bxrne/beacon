#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "button_task.h"
#include "cJSON.h"

extern QueueHandle_t event_queue;
extern TaskHandle_t button_task_handle;
extern TaskHandle_t traffic_light_task_handle;

static void IRAM_ATTR button_isr_handler(void *arg)
{
  BaseType_t xHigherPriorityTaskWoken = pdFALSE;
  vTaskNotifyGiveFromISR(button_task_handle, &xHigherPriorityTaskWoken);
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
  gpio_install_isr_service(0);
  gpio_isr_handler_add(PED_BUTTON_PIN, button_isr_handler, (void *)PED_BUTTON_PIN);

  button_state_t button_state = BUTTON_RELEASED;
  TickType_t last_press_time = 0;

  while (1)
  {
    ulTaskNotifyTake(pdTRUE, portMAX_DELAY);

    TickType_t current_time = xTaskGetTickCount();

    // Handle debouncing
    if (button_state == BUTTON_RELEASED &&
        (current_time - last_press_time) >= pdMS_TO_TICKS(DEBOUNCE_TIME_MS))
    {
      // Verify button is still pressed after debounce delay
      vTaskDelay(pdMS_TO_TICKS(DEBOUNCE_TIME_MS));
      if (gpio_get_level(PED_BUTTON_PIN) == 0)
      {
        button_state = BUTTON_PRESSED;
        last_press_time = current_time;

        ESP_LOGI("EVENT", "Valid button press detected");
        xTaskNotify(traffic_light_task_handle, EVENT_BUTTON_PRESS, eSetValueWithOverwrite);
      }
    }

    // Reset button state when released
    if (gpio_get_level(PED_BUTTON_PIN) == 1)
    {
      button_state = BUTTON_RELEASED;
    }
  }
}
