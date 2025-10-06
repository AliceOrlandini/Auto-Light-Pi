import 'package:auto_light_pi/di/di.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_bloc.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_state.dart';
import 'package:auto_light_pi/features/home_screen.dart';
import 'package:auto_light_pi/navigation/refresh_stream.dart';
import 'package:auto_light_pi/splash.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:auto_light_pi/features/login/presentation/screens/login_screen.dart';

final AuthenticationBloc authBloc = sl<AuthenticationBloc>();

final GoRouter router = GoRouter(
  initialLocation: '/',
  refreshListenable: GoRouterRefreshStream(authBloc.stream),
  redirect: (BuildContext context, GoRouterState state) {
    final AuthenticationState authState = context
        .read<AuthenticationBloc>()
        .state;
    if (authState.status == AuthenticationStatus.unauthenticated) {
      return '/login';
    }
    if (authState.status == AuthenticationStatus.authenticated) {
      return '/home';
    }
    return null; // splash screen
  },
  routes: <GoRoute>[
    GoRoute(
      path: '/',
      name: 'splash',
      builder: (BuildContext context, GoRouterState state) =>
          const SplashScreen(),
    ),
    GoRoute(
      path: '/login',
      name: 'login',
      builder: (BuildContext context, GoRouterState state) =>
          const LoginScreen(),
    ),
    GoRoute(
      path: '/home',
      name: 'home',
      builder: (BuildContext context, GoRouterState state) =>
          const HomeScreen(),
    ),
  ],
);
