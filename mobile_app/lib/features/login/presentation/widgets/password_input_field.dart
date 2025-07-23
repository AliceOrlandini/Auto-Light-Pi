import 'package:auto_light_pi/core/theme/colors.dart';
import 'package:auto_light_pi/core/widgets/generic_input_field.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_event.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_state.dart';
import 'package:auto_light_pi/features/login/presentation/validations/password.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class PasswordInputField extends StatefulWidget {
  const PasswordInputField({super.key});

  @override
  State<PasswordInputField> createState() => _PasswordInputFieldState();
}

class _PasswordInputFieldState extends State<PasswordInputField> {
  bool _obscureText = true;

  void _togglePasswordVisibility() {
    setState(() {
      _obscureText = !_obscureText;
    });
  }

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<LoginBloc, LoginState>(
      buildWhen: (LoginState prev, LoginState curr) =>
          prev.password.value != curr.password.value ||
          prev.isPasswordValid != curr.isPasswordValid,
      builder: (BuildContext context, LoginState state) {
        return GenericInputField(
          label: 'Password',
          keyboardType: TextInputType.text,
          obscureText: _obscureText,
          errorText: () {
            if (state.password.isPure) return null;

            switch (state.password.error) {
              case PasswordValidationError.empty:
              case PasswordValidationError.tooShort:
              case PasswordValidationError.digitMissing:
              case PasswordValidationError.upperCaseMissing:
                return 'Invalid Password';
              default:
                return null;
            }
          }(),
          isValid: state.isPasswordValid,
          suffixIcon: IconButton(
            onPressed: _togglePasswordVisibility,
            icon: Icon(_obscureText ? Icons.visibility : Icons.visibility_off),
            color: AppColors.white,
            iconSize: 20,
          ),
          onChanged: (String value) =>
              context.read<LoginBloc>().add(LoginPasswordChanged(value)),
        );
      },
    );
  }
}
