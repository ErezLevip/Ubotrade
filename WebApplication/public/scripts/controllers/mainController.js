app.controller("mainCtrl", function ($scope, $state, notificationService, $interval, googleOauth) {
    $scope.notifications = [];
    $scope.previusNotifications = [];
    $scope.newNotifications = 0;
    $scope.logoTitle = "UbotTrade";
    $scope.state = "dashboard";
    $scope.go = function (page) {
        $scope.state = page;
        $state.go(page);
    };
    $scope.user = null;
    $scope.notificationWindowOpen = false;

    var init = function () {

        //setup image and first name
        googleOauth.getLoggedinUser().then(function (user) {
            $scope.user = {
                img: user.w3.Paa,
                firstName: user.w3.ofa,
                fullName: user.w3.ig
            };
        });

        $interval(function () {
            getNotifications(false);
        }, 10000);

        getNotifications(false);
    };

    $scope.toggleNotification = function () {
        setTimeout(function () {
            $scope.notificationWindowOpen = $('#notificationIcon').hasClass('open');
            if ($scope.notificationWindowOpen) {
                readAllNotifications(true)
            }
        }, 500);
    };

    var readAllNotifications = function (readAll) {
        notificationService.readAll(readAll);
    };

    var getNotifications = function () {
        notificationService.get().then(function (res) {
            if (res) {
                $scope.previusNotifications = angular.copy($scope.notifications);
                $scope.notifications = [];
                $scope.newNotifications = 0;
                if(res && Array.isArray(res) && res.length > 0) {
                    res.forEach(function (msg) {
                        var message = msg["NotificationMessage"];
                        if (message) {
                            $scope.notifications.push(message);
                            if ($scope.previusNotifications.indexOf(message) == -1)
                                $scope.newNotifications++;
                        }
                    });
                }
            }
        });
    };


    $scope.$on("login", function (event, args) {
        init();
    });
});