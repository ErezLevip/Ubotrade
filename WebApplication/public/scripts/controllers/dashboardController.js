app.controller("dashboardCtrl", function ($scope, $http, $interval, notificationService,googleOauth) {

    $scope.lastProfitsUpdate = new Date();
    $scope.lastTickerUpdate = new Date();

    $scope.botNumber = 0;
    $scope.botName = "";
    $scope.amount = 0;
    $scope.bots = [];
    $scope.botToDisplay = null;
    $scope.activities = [];
    $scope.lastPrice = 0;
    $scope.selectedBot = null;
    $scope.botNotificationCount = 0;
    $scope.showProfitsData = false;
    $scope.isLoggedIn = false;
    /*$scope.$on('login', function (evet, data) {
        $scope.isLoggedIn = true;
    });*/

    googleOauth.isSignedIn().then(function (isSignedIn) {
        if (isSignedIn) {
            googleOauth.getLoggedinUser().then(function (user) {
                $scope.isLoggedIn = true;
            });
        }
    });

    var getBotInfo = function (botNumber) {
        $http.post("/GetBotInformation", {botNumber: botNumber}, null).then(function (res) {
            if (res.data) {
                $scope.amount = res.data["Amount"];
                $scope.botCurrency = res.data["Currency"];
                $scope.botName = res.data["BotName"];
                $scope.botNumber = res.data["BotNumber"];
            }
        }, function (err) {
            console.log(err)
        })
    }

    var getLastActivities = function (botNumber) {
        $http.post("/GetLastActivities", {botNumber: botNumber}, null).then(function (res) {
            if (res.data) {
                $scope.activities = res.data;
            }
        }, function (err) {
            console.log(err)
        })
    }

    $scope.initBotInfo = function (botNumber) {
        var intBotNumber = parseInt(botNumber);
        $scope.botNotificationCount = 0;
        getBotInfo(intBotNumber);
        getLastActivities(intBotNumber);
        initProfitsChartData(intBotNumber);
        initRealTimeChartData(intBotNumber);
        initRealTimeInterval(intBotNumber);
        initBotNotificationCount(intBotNumber)
    };

    var init = function () {
        $scope.$on('login', function () {
            $scope.isLoggedIn = true;
        });

        $http.post("/GetAllActiveBots", {}, null).then(function (res) {
            if (res.data) {
                $scope.bots = res.data;
                if ($scope.bots.length) {
                    $scope.selectedBot = $scope.bots[0];
                    $scope.initBotInfo($scope.selectedBot.BotNumber);
                }
            }
        }, function (err) {
            console.log(err)
        })
    };

    function initBotNotificationCount(botNumber) {
        notificationService.get().then(function (res) {
            if (res) {
                if (res && Array.isArray(res) && res.length >0) {
                    res.forEach(function (msg) {
                        if (msg["BotNumber"]) {
                            if (msg["BotNumber"] == botNumber) {
                                $scope.botNotificationCount++;
                            }
                        }
                    });
                }
            }
        });
    }

    function initProfitsChartData(botNumber) {
        var maxDays = 7;
        $scope.profitsData = [];
        $http.post("/GetBotProfits", {botNumber: botNumber, days: maxDays}, null).then(function (res) {
            $scope.lastProfitsUpdate = new Date();
            if (res.data && res.data.length) {
                $scope.showProfitsData = true;
                $scope.profitsData = res.data;
                var missing = maxDays - $scope.profitsData.length;
                if (missing) {
                    maxDays = 2;
                    $scope.profitsData.push(0.0);
                } else if (maxDays > $scope.profitsData.length) {
                    maxDays = $scope.profitsData.length;
                }

                var dateStrings = [maxDays];
                for (var i = maxDays - 1; i >= 0; i--) {
                    var date = new Date();
                    date.setDate(date.getDate() - i);
                    dateStrings[i] = date.toLocaleDateString();
                }

                var dataEmailsSubscriptionChart = {
                    labels: dateStrings.reverse(),
                    series: [
                        $scope.profitsData
                    ]
                };
                var optionsEmailsSubscriptionChart = {
                    axisX: {
                        showGrid: false
                    },
                    low: 0,
                    high: Math.max.apply(null, $scope.profitsData) * 1.1,
                    chartPadding: {top: 0, right: 5, bottom: 0, left: 0}
                };
                var responsiveOptions = [
                    ['screen and (max-width: 640px)', {
                        seriesBarDistance: 5,
                        axisX: {
                            labelInterpolationFnc: function (value) {
                                return value[0];
                            }
                        }
                    }]
                ];
                var emailsSubscriptionChart = Chartist.Bar('#emailsSubscriptionChart', dataEmailsSubscriptionChart, optionsEmailsSubscriptionChart, responsiveOptions);

                //start animation for the Emails Subscription Chart
                md.startAnimationForBarChart(emailsSubscriptionChart);
            }
        }, function (err) {
            console.log(err)
        })
    }

    function initRealTimeChartData(botNumber) {

        $http.post("/GetBotTickerData", {botNumber: botNumber}, null).then(function (res) {
                $scope.lastTickerUpdate = new Date();
                if (res.data && res.data.length) {
                    $scope.tickerData = res.data;
                    var upperStair = [];
                    var middleStair = [];
                    var lowerStair = [];
                    var priceLine = [];
                    var highLine = [];

                    for (var i = $scope.tickerData.length - 1; i >= 0; i--) {
                        upperStair.push($scope.tickerData[i].Stairs[0]);
                        middleStair.push($scope.tickerData[i].Stairs[1]);
                        lowerStair.push($scope.tickerData[i].Stairs[2]);
                        priceLine.push($scope.tickerData[i].P);
                        highLine.push($scope.tickerData[i].H);
                    }

                    $scope.lastPrice = $scope.tickerData[0].P;

                    var graphData = [];
                    if (middleStair[0] && lowerStair[0]) {
                        graphData = [
                            upperStair, middleStair, lowerStair, priceLine
                        ];
                    } else {
                        graphData = [
                            upperStair, priceLine
                        ];
                        if (highLine[0] && highLine[highLine.length - 1]) {
                            graphData.push(highLine);
                        }
                    }

                    dataCompletedTasksChart = {
                        labels: ['1', '2', '3', '4', '5', '6', '7', '8', '9', '10'],
                        series: graphData
                    };

                    optionsCompletedTasksChart = {
                        lineSmooth: Chartist.Interpolation.cardinal({
                            tension: 0
                        }),
                        low: Math.min(priceLine) * 0.97,
                        high: Math.max.apply(null, upperStair) * 1.02,
                        showPoint: true,
                        chartPadding: {top: 0, right: 0, bottom: 0, left: 0}
                    }

                    var completedTasksChart = new Chartist.Line('#completedTasksChart', dataCompletedTasksChart, optionsCompletedTasksChart);

                    // start animation for the Completed Tasks Chart - Line Chart
                    md.startAnimationForLineChart(completedTasksChart);
                }
            },
            function (err) {
                console.log(err)
            })
    }

    function initRealTimeInterval(botNumber) {
        var stop = $interval(function () {
            var intBotNumber = parseInt(botNumber);
            getBotInfo(intBotNumber);
            getLastActivities(intBotNumber);
            initProfitsChartData(intBotNumber);
            initRealTimeChartData(intBotNumber);
            initRealTimeChartData(intBotNumber);
        }, 60000)
    }

    init();
});