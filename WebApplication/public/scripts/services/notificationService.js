app.factory('notificationService', ['$http', function ($http) {
    var current = [];
    var getAllPromise = null;

    var get = function () {
        return getAllPromise || (getAllPromise = $http.post('/GetNotifications', {}, null).then(function (res) {
            current = res.data;
            return res.data;
        }))
    };

    var readAll = function (readAll) {
        var data = {
            read_all: readAll
        };

        $http.post('/GetNotifications', data, null);
        return null;
    };

    return {
        get: get,
        readAll: readAll
    }

}]);
