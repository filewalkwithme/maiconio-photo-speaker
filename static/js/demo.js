/**
 * Copyright 2014, 2015 IBM Corp. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*global $:false, SPEECH_SYNTHESIS_VOICES */

'use strict';

//this file is a short version of https://github.com/watson-developer-cloud/text-to-speech-nodejs/blob/master/public/js/demo.js
$(document).ready(function() {

  function synthesizeRequest(audio, text) {
    var downloadURL = 'sound?text='+text;

    audio.pause();
    try {
      audio.currentTime = 0;
    } catch(ex) {
      // ignore. Firefox just freaks out here for no apparent reason.
    }
    audio.src = downloadURL;
    audio.play();
    return true;
  }

  var audio = $('.audio').get(0);

  $('.menu-item').click(function(evt) {
    var text = $(this).text();
    text = encodeURIComponent(text)
    synthesizeRequest(audio, text);
  });

  $('.speak-all').click(function(evt) {
    var thisButton = $('#'+this.id)

    var filename = thisButton.parent().find(".file-name").text();

    var labels = thisButton.parent().parent().find(".td-labels").find(".list-labels").children();
    var text = "The photo " + filename + " contains the following labels: "

    labels.each(function() {
      text = text + $(this).find(".menu-item").text()+", ";
    });
    synthesizeRequest(audio, text);
  });
});
