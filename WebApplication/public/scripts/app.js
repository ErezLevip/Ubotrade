var app = angular.module('app', ['ui.router', 'ui.select', 'ngSanitize']);
Window.isLoggedIn = false;
app.config(function ($stateProvider, $urlRouterProvider) {

    $urlRouterProvider.otherwise('/dashboard');
    $stateProvider
        .state('dashboard', {
            url: '/dashboard',
            templateUrl: 'dashboard.html'
        })

        .state('login', {
            url: '/login',
            templateUrl: 'login.html'
        })

        .state('newbot', {
            url: '/newbot',
            templateUrl: 'newbot.html'
        });
});

app.config(['$httpProvider', function ($httpProvider) {
    $httpProvider.interceptors.push('authInterceptor');
}]);

