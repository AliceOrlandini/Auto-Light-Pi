import 'package:flutter/material.dart';
import 'package:auto_light_pi/core/theme/light_theme.dart';
import 'package:auto_light_pi/core/theme/dark_theme.dart';

class AppTheme {
  static ThemeData getTheme(ThemeMode mode) {
    switch (mode) {
      case ThemeMode.system:
      case ThemeMode.light:
        return lightTheme;
      case ThemeMode.dark:
        return darkTheme;
    }
  }
}
