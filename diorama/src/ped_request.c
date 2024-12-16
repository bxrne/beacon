#include <freertos/FreeRTOS.h> 
#include "ped_request.h"
#include "driver/gpio.h"
#include "freertos/queue.h"
#include "config.h"
#include "globals.h"

void IRAM_ATTR button_isr_handler(void *arg)
{
  uint32_t button_pressed = 1;
  xQueueSendFromISR(pedestrianRequestQueue, &button_pressed, NULL);
}

void init_ped_request(void)
{
  esp_rom_gpio_pad_select_gpio(PED_BUTTON_PIN);
  gpio_set_direction(PED_BUTTON_PIN, GPIO_MODE_INPUT);
  gpio_set_intr_type(PED_BUTTON_PIN, GPIO_INTR_POSEDGE);

  gpio_install_isr_service(0);
  gpio_isr_handler_add(PED_BUTTON_PIN, button_isr_handler, NULL);
}
