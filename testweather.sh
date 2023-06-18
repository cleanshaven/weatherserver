#!/bin/bash

dbus-send --dest=com.github.cleanshaven.WeatherService --print-reply \
          --type=method_call \
          /com/github/cleanshaven/WeatherService \
          com.github.cleanshaven.WeatherService.GetWeather \
          string:'-81.86' string:'41.26'

