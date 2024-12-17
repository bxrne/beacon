#ifndef GLOBALS_H
#define GLOBALS_H

#include <freertos/FreeRTOS.h>
#include "freertos/queue.h"
#include "freertos/semphr.h"

extern QueueHandle_t pedestrianRequestQueue;
extern SemaphoreHandle_t xPedestrianSemaphore;

#endif // GLOBALS_H