#ifndef HTTP_CLIENT_H
#define HTTP_CLIENT_H

#include "esp_err.h"

esp_err_t http_post(const char *url, const char *data);

#endif // HTTP_CLIENT_H