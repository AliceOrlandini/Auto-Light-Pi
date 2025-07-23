import 'package:auto_light_pi/core/widgets/generic_input_field.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_event.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_state.dart';
import 'package:auto_light_pi/features/login/presentation/validations/username.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class UsernameInputField extends StatelessWidget {
  const UsernameInputField({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<LoginBloc, LoginState>(
      buildWhen: (LoginState prev, LoginState curr) =>
          prev.username.value != curr.username.value ||
          prev.isUsernameValid != curr.isUsernameValid,
      builder: (BuildContext context, LoginState state) {
        return GenericInputField(
          label: 'Username',
          keyboardType: TextInputType.text,
          errorText: () {
            if (state.username.isPure) return null;

            return state.username.error == UsernameValidationError.empty
                ? 'Please, insert the username'
                : null;
          }(),
          isValid: state.isUsernameValid,
          onChanged: (String value) =>
              context.read<LoginBloc>().add(LoginUsernameChanged(value)),
        );
      },
    );
  }
}
