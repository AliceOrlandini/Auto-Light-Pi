import 'package:auto_light_pi/di/di.dart';
import 'package:flutter/material.dart';
import 'app.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  init();
  runApp(const App());
}
