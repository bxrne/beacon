#ifndef METRICS_H
#define METRICS_H

#include <stddef.h> // Include for size_t

void get_metrics(char *buffer, size_t buffer_size);
void init_metrics_buffers(size_t size);
void update_light_buffers();
void free_metrics_buffers();

#endif // METRICS_H
