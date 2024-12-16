#include <freertos/FreeRTOS.h>
#include "ped_light.h"
#include "freertos/queue.h"
#include "esp_log.h"
#include "globals.h"

#define TAG "PED_LIGHT"

void PedestrianLightTask(void *pvParameters)
{
  uint32_t button_pressed;
  while (true)
  {
    if (xQueueReceive(pedestrianRequestQueue, &button_pressed, portMAX_DELAY))
    {
      ESP_LOGI(TAG, "Pedestrian button pressed");
      // Signal the car light task
      xSemaphoreGive(xPedestrianSemaphore);
    }
  }
}
