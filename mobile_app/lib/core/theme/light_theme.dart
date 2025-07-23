import 'package:auto_light_pi/core/theme/colors.dart';
import 'package:flutter/material.dart';

final ThemeData lightTheme = ThemeData(
  colorScheme: ColorScheme.fromSwatch(
    brightness: Brightness.light,
    primarySwatch: AppColors.primary,
  ),
  fontFamily: 'Montserrat',
  useMaterial3: true,
);
