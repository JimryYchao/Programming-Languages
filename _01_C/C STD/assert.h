#pragma once

// #define NDEBUG
#ifndef NDEBUG
void assert(scalar_expression);            // Assertion failed: expression, function abc, file xyz, line nnn.
#else
void assert(...)((void)0);
#endif