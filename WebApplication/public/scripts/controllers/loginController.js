app.controller("loginCtrl", function ($scope, $state, googleOauth, $injector, $rootScope) {
    $scope.login = function () {
        googleOauth.redirectToLogin();
        invokeLogin();
    };
    $scope.logoTitle = "UbotTrade";
    var invokeLogin = function () {
        googleOauth.isSignedIn().then(function (isSignedIn) {
            if (isSignedIn) {
                googleOauth.getLoggedinUser().then(function (user) {
                    var $http = $injector.get('$http');
                    $http.post("/Login", {
                        session_id: user.Zi.access_token,
                        uid: user.El,
                        first_name: user.w3.ofa,
                        last_name: user.w3.Wea,
                        email: user.w3.U3
                    }, {loginRequest: true}).then(function (res) {
                        if (res.data.IsAuthorized) {
                            setTimeout(function () {
                                $rootScope.$broadcast('login')
                            }, 1000);
                            $('#lean_overlay').remove();
                            $state.go('dashboard');
                        } else {
                            alert("error check console");
                            console.log(res.data.IsAuthorized);
                        }
                    })
                });
            }
            else {
                $("#modal_trigger").leanModal({
                    top: 100,
                    overlay: 0.6,
                    closeButton: ".modal_close"
                });
                $("#modal_trigger").click();
            }
        });
    };
    setTimeout(function () {
    invokeLogin();
    },500);
});
