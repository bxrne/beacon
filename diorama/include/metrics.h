#ifndef METRICS_H
#define METRICS_H

#include <stddef.h>
#include "freertos/FreeRTOS.h"
#include "freertos/semphr.h"

typedef enum
{
  LIGHT_GREEN,
  LIGHT_YELLOW,
  LIGHT_RED
} LightColor;

void init_metrics_buffers(size_t size);

void add_car_light_state(LightColor state);
void add_ped_light_state(LightColor state);

LightColor get_recent_car_light_state();
LightColor get_recent_ped_light_state();

const char *light_color_to_string(LightColor color);

void free_metrics_buffers();

void get_current_time_utc(char *buffer, size_t buffer_size);

#endif // METRICS_H
