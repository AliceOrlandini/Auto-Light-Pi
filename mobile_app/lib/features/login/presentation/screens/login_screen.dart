import 'package:auto_light_pi/core/widgets/error_toast.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_state.dart';
import 'package:auto_light_pi/features/login/presentation/widgets/login_button.dart';
import 'package:flutter/material.dart';
import 'package:auto_light_pi/di/di.dart' as di;
import 'package:auto_light_pi/features/login/presentation/widgets/username_input_field.dart';
import 'package:auto_light_pi/features/login/presentation/widgets/password_input_field.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:formz/formz.dart';
import 'package:go_router/go_router.dart';

class LoginScreen extends StatelessWidget {
  const LoginScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocProvider<LoginBloc>(
      create: (_) => di.sl<LoginBloc>(),
      child: Scaffold(
        appBar: AppBar(
          centerTitle: true,
          title: Text(
            'Auto Light Pi'.toUpperCase(),
            style: const TextStyle(fontWeight: FontWeight.bold),
          ),
        ),
        body: BlocListener<LoginBloc, LoginState>(
          listenWhen: (LoginState previous, LoginState current) {
            // listen only when the status passess from non-failure to failure
            // or non-success to success
            return (previous.status != FormzSubmissionStatus.failure &&
                    current.status == FormzSubmissionStatus.failure) ||
                (previous.status != FormzSubmissionStatus.success &&
                    current.status == FormzSubmissionStatus.success) ||
                (previous.status == FormzSubmissionStatus.initial &&
                    (current.status == FormzSubmissionStatus.success ||
                        current.status == FormzSubmissionStatus.failure));
          },
          listener: (BuildContext context, LoginState state) {
            switch (state.status) {
              case FormzSubmissionStatus.success:
                context.go('/home');
                break;
              case FormzSubmissionStatus.failure:
                ErrorToast.show(context, description: state.errorMessage);
                break;
              default:
                break;
            }
          },
          child: const Padding(
            padding: EdgeInsets.all(16.0),
            child: Column(
              children: <Widget>[
                SizedBox(height: 50),
                UsernameInputField(),
                SizedBox(height: 30),
                PasswordInputField(),
                SizedBox(height: 30),
                LoginButton(),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
