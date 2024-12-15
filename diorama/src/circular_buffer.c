#include "circular_buffer.h"
#include <stdlib.h>
#include <string.h>

CircularBuffer *circular_buffer_init(size_t size)
{
  CircularBuffer *cb = malloc(sizeof(CircularBuffer));
  if (cb == NULL)
  {
    return NULL;
  }
  cb->buffer = malloc(size);
  if (cb->buffer == NULL)
  {
    free(cb);
    return NULL;
  }
  cb->max = size;
  circular_buffer_reset(cb);
  return cb;
}

void circular_buffer_free(CircularBuffer *cb)
{
  free(cb->buffer);
  free(cb);
}

void circular_buffer_reset(CircularBuffer *cb)
{
  cb->head = 0;
  cb->tail = 0;
  cb->full = 0;
}

size_t circular_buffer_write(CircularBuffer *cb, const char *data, size_t bytes)
{
  size_t capacity = cb->max;
  size_t bytes_to_write = bytes;

  if (cb->full)
  {
    return 0;
  }

  if (bytes > capacity - cb->head)
  {
    bytes_to_write = capacity - cb->head;
  }

  memcpy(cb->buffer + cb->head, data, bytes_to_write);
  cb->head = (cb->head + bytes_to_write) % capacity;

  if (cb->head == cb->tail)
  {
    cb->full = 1;
  }

  return bytes_to_write;
}

size_t circular_buffer_read(CircularBuffer *cb, char *data, size_t bytes)
{
  size_t capacity = cb->max;
  size_t bytes_to_read = bytes;

  if (cb->head == cb->tail && !cb->full)
  {
    return 0;
  }

  if (bytes > capacity - cb->tail)
  {
    bytes_to_read = capacity - cb->tail;
  }

  memcpy(data, cb->buffer + cb->tail, bytes_to_read);
  cb->tail = (cb->tail + bytes_to_read) % capacity;
  cb->full = 0;

  return bytes_to_read;
}
