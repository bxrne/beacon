#include "metrics.h"
#include "circular_buffer.h"
#include "freertos/FreeRTOS.h"
#include "freertos/semphr.h"
#include "esp_log.h"
#include "globals.h"
#include "config.h"
#include "driver/gpio.h" // Include for gpio_get_level
#include <time.h>
#include <sys/time.h>
#include <stdlib.h>

extern SemaphoreHandle_t xPedestrianSemaphore;

// Replace CircularBuffer with circular_buffer_t
static circular_buffer_t *car_light_buffer = NULL;
static circular_buffer_t *ped_light_buffer = NULL;

void init_metrics_buffers(size_t size)
{
  car_light_buffer = circular_buffer_init(size);
  ped_light_buffer = circular_buffer_init(size);
}

void add_car_light_state(LightColor state)
{
  if (car_light_buffer)
  {
    circular_buffer_push(car_light_buffer, (int)state);
  }
}

void add_ped_light_state(LightColor state)
{
  if (ped_light_buffer)
  {
    circular_buffer_push(ped_light_buffer, (int)state);
  }
}

LightColor get_recent_car_light_state()
{
  if (car_light_buffer && !circular_buffer_is_empty(car_light_buffer))
  {
    return (LightColor)circular_buffer_peek_last(car_light_buffer);
  }
  return LIGHT_RED; // Default value if buffer is empty
}

LightColor get_recent_ped_light_state()
{
  if (ped_light_buffer && !circular_buffer_is_empty(ped_light_buffer))
  {
    return (LightColor)circular_buffer_peek_last(ped_light_buffer);
  }
  return LIGHT_RED; // Default value if buffer is empty
}

const char *light_color_to_string(LightColor color)
{
  switch (color)
  {
  case LIGHT_GREEN:
    return "GREEN";
  case LIGHT_YELLOW:
    return "YELLOW";
  case LIGHT_RED:
    return "RED";
  default:
    return "UNKNOWN";
  }
}

void free_metrics_buffers()
{
  if (car_light_buffer)
  {
    circular_buffer_free(car_light_buffer);
    car_light_buffer = NULL;
  }
  if (ped_light_buffer)
  {
    circular_buffer_free(ped_light_buffer);
    ped_light_buffer = NULL;
  }
}

// Remove the update_light_buffers and get_metrics functions if they're no longer needed
