#ifndef CIRCULAR_BUFFER_H
#define CIRCULAR_BUFFER_H

#include <stddef.h>

typedef struct
{
  char *buffer;
  size_t head;
  size_t tail;
  size_t max;
  int full;
} CircularBuffer;

CircularBuffer *circular_buffer_init(size_t size);
void circular_buffer_free(CircularBuffer *cb);
void circular_buffer_reset(CircularBuffer *cb);
size_t circular_buffer_write(CircularBuffer *cb, const char *data, size_t bytes);
size_t circular_buffer_read(CircularBuffer *cb, char *data, size_t bytes);

#endif // CIRCULAR_BUFFER_H
