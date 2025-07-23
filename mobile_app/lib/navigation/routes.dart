import 'package:auto_light_pi/features/home_screen.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:auto_light_pi/features/login/presentation/screens/login_screen.dart';

final GoRouter router = GoRouter(
  initialLocation: '/login',
  routes: <GoRoute>[
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
