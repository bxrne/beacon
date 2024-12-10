#ifndef WIFI_HANDLER_H
#define WIFI_HANDLER_H

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include "esp_event.h" // For esp_event_base_t

void wifi_init_sta(void);
void wait_for_wifi_connection(void);

#endif // WIFI_HANDLER_H
