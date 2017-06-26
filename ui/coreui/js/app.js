/*****
* CONFIGURATION
*/

  // In-browser Storage
  var ibstorage = window.localStorage;
  var default_ep = '';
  var default_remote = '';

  // ibstorage.setItem(key, val);
  // ibstorage.getItem(key);
  // ibstorage.removeItem(key);

  //Main navigation
  $.navigation = $('nav > ul.nav');

  $.panelIconOpened = 'icon-arrow-up';
  $.panelIconClosed = 'icon-arrow-down';

  //Default colours
  $.brandPrimary =  '#20a8d8';
  $.brandSuccess =  '#4dbd74';
  $.brandInfo =     '#63c2de';
  $.brandWarning =  '#f8cb00';
  $.brandDanger =   '#f86c6b';

  $.grayDark =      '#2a2c36';
  $.gray =          '#55595c';
  $.grayLight =     '#818a91';
  $.grayLighter =   '#d1d4d7';
  $.grayLightest =  '#f8f9fa';

'use strict';

/****
* MAIN NAVIGATION
*/

$(document).ready(function($){

  // Check if we have any defaults:
  default_ep = ibstorage.getItem('reshifter.info/default-etcd');
  if (default_ep !== '') {
    console.info('reshifter.info/default-etcd:'+default_ep)
    $('#endpoint').val(default_ep);
  }
  default_remote = ibstorage.getItem('reshifter.info/default-remote');
  if (default_remote !== '') {
    console.info('reshifter.info/default-remote:'+default_remote)
    $('#backup-result').html('<div>Backing up to: <code>'+ default_remote +'</code></div>');
  }

  // ACTIONS:
  $('#doexplore').click(function(event) {
    var ep = $('#endpoint').val();
    $.ajax({
        type: "GET",
        url: 'http://localhost:8080/v1/explorer?endpoint='+encodeURIComponent(ep),
        dataType: 'json',
        async: false,
        data: '{"endpoint": "' + ep +'"}',
        error: function (d) {
          console.info(d);
          $('#config-result').html('<h2>Result</h2>')
          $('#config-result').append('<div>There was a problem carrying out the config:<br><code>'+ d.responseText + '</code> </div>')
        },
        success: function (d) {
          console.info(d);
          $('#config-result').html('<h2>Result</h2>')
            $('#config-result').append('<div>etcd: <code>v' + d.etcdversion +', ' + d.etcdsec +'</code></div>')
            $('#config-result').append('<div>Kubernetes: <code>' + d.k8sdistro +'</code></div>')
        }
    })
  });

  $('#dosaveconfig').click(function(event) {
    var endpoint = $('#endpoint').val();
    var remote = $('#remote:checked').val();
    if (remote === 's3'){
      remote += ':'+ $('#bucket').val();
    }
    ibstorage.setItem('reshifter.info/default-etcd', endpoint);
    console.info('reshifter.info/default-etcd:'+endpoint)
    ibstorage.setItem('reshifter.info/default-remote', remote);
    console.info('reshifter.info/default-remote:'+remote)
    $('#config-result').html('<h2>Result</h2><div>All settings stored locally.</div>')
  });

  $('#dobackup').click(function(event) {
    var ep = $('#endpoint').val();
    $.ajax({
        type: "POST",
        url: 'http://localhost:8080/v1/backup',
        dataType: 'json',
        async: false,
        data: '{"endpoint": "' + ep +'", "bucket": "' + default_remote +'" }',
        error: function (d) {
          console.info(d);
          $('#backup-result').html('<h2>Result</h2>')
          $('#backup-result').append('<div>There was a problem carrying out the backup:<br><code>'+ d.responseText + '</code> </div>')
        },
        success: function (d) {
          console.info(d);
          $('#backup-result').html('<h2>Result</h2>')
          if(d.outcome == 'success'){
            $('#backup-result').append('<div>The backup with ID <code>' + d.backupid +'</code> is now available <a href="/v1/backup/'+ d.backupid + '">here</a> for download.</div>')
          } else{
            $('#backup-result').append('<div>There was a problem carrying out the backup:<br><pre>'+ d + '"</pre> </div>')
          }
        }
    })
  });

  $('#dorestore').click(function(event) {
    var ep = $('#endpoint').val();
    var bid = $('#backupid').val();
    $.ajax({
        type: "POST",
        url: 'http://localhost:8080/v1/restore',
        dataType: 'json',
        async: false,
        data: '{ "endpoint": "' + ep +'", "archive": "' + bid +'" }',
        error: function (d) {
          console.info(d);
          $('#restore-result').html('<h2>Result</h2>')
          $('#restore-result').append('<div>There was a problem carrying out the restore:<br><code>'+ d.responseText + '</code> </div>')
        },
        success: function (d) {
          console.info(d);
          $('#restore-result').html('<h2>Result</h2>')
          if(d.outcome == 'success'){
            $('#restore-result').append('<div>Restored ' + d.keysrestored + ' keys from backup with ID <code>' + bid +'</code> to <code>'+ ep + '</code>.</div>')
          } else{
            $('#restore-result').append('<div>There was a problem carrying out the restore:<br><pre>'+ d + '</pre> </div>')
          }
        }
    })
  });


  // Add class .active to current link
  $.navigation.find('a').each(function(){

    var cUrl = String(window.location).split('?')[0];

    if (cUrl.substr(cUrl.length - 1) == '#') {
      cUrl = cUrl.slice(0,-1);
    }

    if ($($(this))[0].href==cUrl) {
      $(this).addClass('active');

      $(this).parents('ul').add(this).each(function(){
        $(this).parent().addClass('open');
      });
    }
  });

  // Dropdown Menu
  $.navigation.on('click', 'a', function(e){

    if ($.ajaxLoad) {
      e.preventDefault();
    }

    if ($(this).hasClass('nav-dropdown-toggle')) {
      $(this).parent().toggleClass('open');
      resizeBroadcast();
    }

  });


  function resizeBroadcast() {

    var timesRun = 0;
    var interval = setInterval(function(){
      timesRun += 1;
      if(timesRun === 5){
        clearInterval(interval);
      }
      window.dispatchEvent(new Event('resize'));
    }, 62.5);
  }

  /* ---------- Main Menu Open/Close, Min/Full ---------- */
  $('.navbar-toggler').click(function(){

    if ($(this).hasClass('sidebar-toggler')) {
      $('body').toggleClass('sidebar-hidden');
      resizeBroadcast();
    }

    if ($(this).hasClass('sidebar-minimizer')) {
      $('body').toggleClass('sidebar-minimized');
      resizeBroadcast();
    }

    if ($(this).hasClass('aside-menu-toggler')) {
      $('body').toggleClass('aside-menu-hidden');
      resizeBroadcast();
    }

    if ($(this).hasClass('mobile-sidebar-toggler')) {
      $('body').toggleClass('sidebar-mobile-show');
      resizeBroadcast();
    }

  });

  $('.sidebar-close').click(function(){
    $('body').toggleClass('sidebar-opened').parent().toggleClass('sidebar-opened');
  });

  /* ---------- Disable moving to top ---------- */
  $('a[href="#"][data-top!=true]').click(function(e){
    e.preventDefault();
  });

});

/****
* CARDS ACTIONS
*/

$(document).on('click', '.card-actions a', function(e){
  e.preventDefault();

  if ($(this).hasClass('btn-close')) {
    $(this).parent().parent().parent().fadeOut();
  } else if ($(this).hasClass('btn-minimize')) {
    var $target = $(this).parent().parent().next('.card-block');
    if (!$(this).hasClass('collapsed')) {
      $('i',$(this)).removeClass($.panelIconOpened).addClass($.panelIconClosed);
    } else {
      $('i',$(this)).removeClass($.panelIconClosed).addClass($.panelIconOpened);
    }

  } else if ($(this).hasClass('btn-setting')) {
    $('#myModal').modal('show');
  }

});

function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

function init(url) {

  /* ---------- Tooltip ---------- */
  $('[rel="tooltip"],[data-rel="tooltip"]').tooltip({"placement":"bottom",delay: { show: 400, hide: 200 }});

  /* ---------- Popover ---------- */
  $('[rel="popover"],[data-rel="popover"],[data-toggle="popover"]').popover();

}
