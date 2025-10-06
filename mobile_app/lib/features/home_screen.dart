import 'package:auto_light_pi/core/theme/colors.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_bloc.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_event.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_state.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final String username = context.select<AuthenticationBloc, String>(
      (AuthenticationBloc bloc) => bloc.state.user.username,
    );

    return Scaffold(
      appBar: AppBar(
        centerTitle: true,
        title: Text(
          'Auto Light Pi'.toUpperCase(),
          style: const TextStyle(fontWeight: FontWeight.bold),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: <Widget>[
            Text(
              'Welcome, $username!',
              style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w500),
            ),
            const SizedBox(height: 20),
            // go to the bottom of the screen
            Expanded(child: Container()),
            ElevatedButton(
              onPressed: () {
                context.read<AuthenticationBloc>().add(LogoutRequest());
                context.go('/');
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.red,
                padding: const EdgeInsets.symmetric(
                  horizontal: 32,
                  vertical: 12,
                ),
              ),
              child:
                  context.read<AuthenticationBloc>().state.status ==
                      AuthenticationStatus.unknown
                  ? const SizedBox(
                      width: 20,
                      height: 20,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        valueColor: AlwaysStoppedAnimation<Color>(
                          AppColors.white,
                        ),
                      ),
                    )
                  : const Text(
                      'Logout',
                      style: TextStyle(fontSize: 16, color: AppColors.white),
                    ),
            ),
          ],
        ),
      ),
    );
  }
}
