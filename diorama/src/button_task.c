#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "button_task.h"
#include "cJSON.h"

extern QueueHandle_t event_queue;
static QueueHandle_t button_evt_queue = NULL;

static void IRAM_ATTR button_isr_handler(void *arg)
{
  uint32_t gpio_num = (uint32_t)arg;
  xQueueSendFromISR(button_evt_queue, &gpio_num, NULL);
}

void log_button_event()
{
  char telemetry[64];
  snprintf(telemetry, sizeof(telemetry), "{\"event\":\"BUTTON_PRESS\"}");
  ESP_LOGI(TAG, "[TELEMETRY] %s", telemetry);
}

void button_task(void *pvParameters)
{
  ESP_LOGI(TAG, "[DEBUG] Button task started");

  // Create queue for button events
  button_evt_queue = xQueueCreate(BUTTON_QUEUE_SIZE, sizeof(uint32_t));
  if (button_evt_queue == NULL)
  {
    ESP_LOGE(TAG, "[ERROR] Failed to create button queue");
    vTaskDelete(NULL);
    return;
  }

  // Configure button GPIO with interrupt
  gpio_config_t io_conf = {
      .pin_bit_mask = (1ULL << PED_BUTTON_PIN),
      .mode = GPIO_MODE_INPUT,
      .pull_up_en = GPIO_PULLUP_ENABLE,
      .pull_down_en = GPIO_PULLDOWN_DISABLE,
      .intr_type = GPIO_INTERRUPT_TRIGGER};

  gpio_config(&io_conf);
  gpio_install_isr_service(0);
  gpio_isr_handler_add(PED_BUTTON_PIN, button_isr_handler, (void *)PED_BUTTON_PIN);

  uint32_t gpio_num;
  button_state_t button_state = BUTTON_RELEASED;
  TickType_t last_press_time = 0;

  while (1)
  {
    if (xQueueReceive(button_evt_queue, &gpio_num, portMAX_DELAY))
    {
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

          ESP_LOGI(TAG, "[EVENT] Valid button press detected");
          event_t event = EVENT_BUTTON_PRESS;
          if (xQueueSend(event_queue, &event, 0) != pdTRUE)
          {
            ESP_LOGE(TAG, "[ERROR] Failed to send button event");
          }
          log_button_event();
        }
      }

      // Reset button state when released
      if (gpio_get_level(PED_BUTTON_PIN) == 1)
      {
        button_state = BUTTON_RELEASED;
      }
    }
  }
}
