#include "metrics.h"
#include "circular_buffer.h"
#include "freertos/FreeRTOS.h"
#include "freertos/semphr.h"
#include "esp_log.h"
#include "globals.h"
#include "config.h"
#include "driver/gpio.h"
#include <time.h>
#include <sys/time.h>
#include <stdlib.h>

extern SemaphoreHandle_t xPedestrianSemaphore;

// Replace CircularBuffer with circular_buffer_t
static circular_buffer_t *car_light_buffer = NULL;
static circular_buffer_t *ped_light_buffer = NULL;

// Declare mutexes for thread safety
static SemaphoreHandle_t car_light_mutex = NULL;
static SemaphoreHandle_t ped_light_mutex = NULL;

void init_metrics_buffers(size_t size)
{
  car_light_buffer = circular_buffer_init(size);
  ped_light_buffer = circular_buffer_init(size);

  // Initialize mutexes
  car_light_mutex = xSemaphoreCreateMutex();
  ped_light_mutex = xSemaphoreCreateMutex();
}

void add_car_light_state(LightColor state)
{
  if (car_light_buffer)
  {
    // Take mutex before accessing the buffer
    xSemaphoreTake(car_light_mutex, portMAX_DELAY);
    circular_buffer_push(car_light_buffer, (int)state);
    xSemaphoreGive(car_light_mutex);
  }
}

void add_ped_light_state(LightColor state)
{
  if (ped_light_buffer)
  {
    xSemaphoreTake(ped_light_mutex, portMAX_DELAY);
    circular_buffer_push(ped_light_buffer, (int)state);
    xSemaphoreGive(ped_light_mutex);
  }
}

LightColor get_recent_car_light_state()
{
  LightColor state = LIGHT_RED; // Default value
  if (car_light_buffer)
  {
    xSemaphoreTake(car_light_mutex, portMAX_DELAY);
    if (!circular_buffer_is_empty(car_light_buffer))
    {
      state = (LightColor)circular_buffer_peek_last(car_light_buffer);
    }
    xSemaphoreGive(car_light_mutex);
  }
  return state;
}

LightColor get_recent_ped_light_state()
{
  LightColor state = LIGHT_RED; // Default value
  if (ped_light_buffer)
  {
    xSemaphoreTake(ped_light_mutex, portMAX_DELAY);
    if (!circular_buffer_is_empty(ped_light_buffer))
    {
      state = (LightColor)circular_buffer_peek_last(ped_light_buffer);
    }
    xSemaphoreGive(ped_light_mutex);
  }
  return state;
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

  // Delete mutexes
  if (car_light_mutex)
  {
    vSemaphoreDelete(car_light_mutex);
    car_light_mutex = NULL;
  }
  if (ped_light_mutex)
  {
    vSemaphoreDelete(ped_light_mutex);
    ped_light_mutex = NULL;
  }
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

// Remove the update_light_buffers and get_metrics functions if they're no longer needed
