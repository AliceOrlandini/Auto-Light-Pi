import 'package:flutter/material.dart';
import 'package:toastification/toastification.dart';

class ErrorToast {
  static void show(BuildContext context, {String? description}) {
    toastification.show(
      context: context,
      title: const Text('Login Failed'),
      description: Text(description ?? 'an error has occurred'),
      type: ToastificationType.error,
      style: ToastificationStyle.fillColored,
      alignment: Alignment.bottomCenter,
      autoCloseDuration: const Duration(seconds: 5),
      animationDuration: const Duration(milliseconds: 300),
      icon: const Icon(Icons.error_outline),
      showIcon: true,
    );
  }
}
