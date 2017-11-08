app.factory('googleOauth', ['$q','$state', function ($q,$state) {
    var GoogleAuth;
    var SCOPE = 'https://www.googleapis.com/auth/drive.metadata.readonly';
    var user = null;
    var deferred = $q.defer();

    (function () {
        var p = document.createElement('script');
        p.type = 'text/javascript';
        p.async = false;
        p.src = 'https://apis.google.com/js/client.js?onload=onGapiLoad';
        var s = document.getElementsByTagName('script')[0];
        s.parentNode.insertBefore(p, s);
    })();

    window.onGapiLoad = function onGapiLoad() {
        gapi.load('client:auth2', init);

    };

    function init() {
        var discoveryUrl = 'https://www.googleapis.com/discovery/v1/apis/drive/v3/rest';
        var clientId = '1016427666052-8dkpso3l7t3aftaon7n9e2hpjq3r8ruk.apps.googleusercontent.com';
        var apiKey = 'AIzaSyBQQi5GAiPfu8CK9ORZRZmlE8c_7N_68VM';
        var SCOPE = 'https://www.googleapis.com/auth/drive.metadata.readonly';
        gapi.client.init({
            'apiKey': apiKey,
            'discoveryDocs': [discoveryUrl],
            'clientId': clientId,
            'scope': SCOPE
        });
        GoogleAuth = gapi.auth2.getAuthInstance();
        if (GoogleAuth) {
            GoogleAuth.isSignedIn.listen(checkIsAuthorized);
            deferred.resolve(null);
        } else {
            deferred.reject(null);
        }
    }

    function checkIsAuthorized(isSignedIn) {
        user = GoogleAuth.currentUser.get();
        var isAuthorized = user.hasGrantedScopes(SCOPE);
        if (!isAuthorized || !isSignedIn) {
            GoogleAuth.signIn().then(function (res) {
                location.reload();
            })
        }
    }

    var isSignedIn = function () {
        var isSignedDefer = $q.defer();
        deferred.promise.then(function (g) {
            isSignedDefer.resolve(GoogleAuth.isSignedIn.get());
        });
        return isSignedDefer.promise;
    };

    var getLoggedinUser = function () {
        var oathDeferred = $q.defer();
        deferred.promise.then(function (g) {
            var signedin = GoogleAuth.isSignedIn.get();
            if (signedin) {
                var user = GoogleAuth.currentUser.get();
                oathDeferred.resolve(user);
            } else {
                $state.go("login");
                oathDeferred.reject(null);
            }
        });
        return oathDeferred.promise;
    };

    var redirectToLogin = function () {
        var oathDeferred = $q.defer();
        deferred.promise.then(function (g) {
            GoogleAuth.signIn().then(function (res) {
                location.reload();
            })
        });
        return oathDeferred.promise;
    };
    return {
        getLoggedinUser: getLoggedinUser,
        isSignedIn: isSignedIn,
        redirectToLogin: redirectToLogin
    }
}]);
