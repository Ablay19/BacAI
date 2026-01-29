#include "text_utils.h"

/**
 * Fast text sanitization in C to avoid Go string allocations.
 * Removes non-printable characters and trims whitespace.
 */
void sanitize_text_in_place(char* s) {
    if (!s) return;
    
    char* src = s;
    char* dst = s;
    
    // Trim leading whitespace
    while (*src && isspace((unsigned char)*src)) src++;
    
    while (*src) {
        if (isprint((unsigned char)*src) || *src == '\n' || *src == '\r' || *src == '\t') {
            *dst++ = *src;
        }
        src++;
    }
    
    *dst = '\0';
    
    // Trim trailing whitespace
    if (dst > s) {
        dst--;
        while (dst >= s && isspace((unsigned char)*dst)) {
            *dst-- = '\0';
        }
    }
}
