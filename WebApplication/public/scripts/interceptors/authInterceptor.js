app.factory('authInterceptor', ['$q', 'googleOauth', '$state', function ($q, googleOauth, $state) {
    return {
        request: function (config) {
            //workaround selectize is using a relative path (open issue)
            if (!checkIgnoreLst(config.url)) {
                if (config.url[0] != '/') {
                    config.url = "/" + config.url;
                }
            }
            console.log(config.url);
            if (config.url.indexOf('.html') != -1 || config.loginRequest) {
                return config || $q.when(config);
            }

            var promise = googleOauth.getLoggedinUser().then(function (user) {
                if (user) {
                    config.headers['X-AUTHORIZATION'] = user.Zi.access_token;
                }
                return config;
            });
            return promise;
        }
    }
}]);


var ignoreList = ["selectize"];

var checkIgnoreLst = function (url) {
    for (var i = 0; i < ignoreList.length; i++) {
        if (url.indexOf(ignoreList[i]) > -1) {
            return true;
        }
    }
};

