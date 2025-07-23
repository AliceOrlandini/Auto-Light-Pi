import 'package:flutter/material.dart';
import 'package:auto_light_pi/core/theme/colors.dart';

class GenericInputField extends StatelessWidget {
  final String label;
  final bool autocorrect;
  final bool enableSuggestions;
  final bool obscureText;
  final TextInputType keyboardType;
  final bool? isValid;
  final String? errorText;
  final ValueChanged<String>? onChanged;
  final Widget? suffixIcon;

  const GenericInputField({
    super.key,
    required this.label,
    this.autocorrect = false,
    this.enableSuggestions = false,
    this.obscureText = false,
    this.keyboardType = TextInputType.text,
    this.isValid,
    this.errorText,
    this.onChanged,
    this.suffixIcon,
  });

  @override
  Widget build(BuildContext context) {
    return TextFormField(
      onChanged: onChanged,
      keyboardType: keyboardType,
      autocorrect: autocorrect,
      enableSuggestions: enableSuggestions,
      obscureText: obscureText,
      decoration: InputDecoration(
        labelText: label,
        errorText: errorText,
        suffixIcon: suffixIcon,
        errorStyle: const TextStyle(color: AppColors.invalidFormField),
        errorBorder: const OutlineInputBorder(
          borderSide: BorderSide(color: AppColors.invalidFormField),
        ),
        focusedBorder: const OutlineInputBorder(
          borderSide: BorderSide(color: AppColors.primary),
        ),
        focusedErrorBorder: const OutlineInputBorder(
          borderSide: BorderSide(color: AppColors.invalidFormField),
        ),
        enabledBorder: OutlineInputBorder(
          borderSide: BorderSide(
            color: isValid == true
                ? AppColors.validFormField
                : AppColors.primary,
          ),
        ),
      ),
      onTapOutside: (PointerDownEvent event) {
        FocusScope.of(context).unfocus();
      },
    );
  }
}
