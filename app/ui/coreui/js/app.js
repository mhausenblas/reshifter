// In-browser Storage
var ibstorage = window.localStorage;

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

$(document).ready(function($){
  getVersion()
  setDefaults()
  initBackup();
  initRestore();

  // ACTIONS:
  $('#doexplore').click(function(event) {
    var ep = $('#endpoint').val();
    $('#config-result').html('<div><img src="./img/standby.gif" alt="please wait" width="64px"></div>');
    $.ajax({
        type: "GET",
        url: 'http://localhost:8080/v1/explorer?endpoint='+encodeURIComponent(ep),
        dataType: 'json',
        async: false,
        data: '{"endpoint": "' + ep +'"}',
        error: function (d) {
          console.info(d);
          $('#config-result').html('<h2>Result</h2>');
          $('#config-result').append('<div>There was a problem exploring the endpoint:<div><code>' + d.responseText + '</code></div></div>');
        },
        success: function (d) {
          console.info(d);
          $('#config-result').html('<h2>Result</h2>');
          $('#config-result').append('<div>etcd: <code>v' + d.etcdversion +', ' + d.etcdsec +'</code></div>');
          $('#config-result').append('<div>distribution: <code>' + d.k8sdistro +'</code></div>');
        }
    })
    $.ajax({
        type: "GET",
        url: 'http://localhost:8080/v1/epstats?endpoint='+encodeURIComponent(ep),
        dataType: 'json',
        async: false,
        data: '{"endpoint": "' + ep +'"}',
        error: function (d) {
          console.info(d);
          $('#config-result').html('<h2>Result</h2>');
          $('#config-result').append('<div>There was a problem collecting endpoint stats:<br><code>'+ d.responseText + '</code> </div>');
          $('#config-result').append('<div style="margin-bottom: 50px"></div>');
        },
        success: function (d) {
          console.info(d);
          $('#config-result').append('<div style="margin-top: 20px">Stats for known keys as per Kubernetes distribution:</div>');
          $('#config-result').append('<div>number of keys: <code>' + d.numkeys +'</code></div>');
          $('#config-result').append('<div>total number of bytes to back up: <code>' + d.totalsizevalbytes +'</code></div>');
          $('#config-result').append('<div style="margin-bottom: 50px"></div>');
        }
    })
  });

  $('#dosaveconfig').click(function(event) {
    var endpoint = $('#endpoint').val();
    var apiversion = $('#apiversion').val();
    var remoteep = $('#remoteep').val();
    var remote = $('#remote:checked').val();
    if (remote === 's3'){
      remote += ' ' + remoteep + ' '+ $('#bucket').val();
    }
    ibstorage.setItem('reshifter.info/default-etcd', endpoint);
    console.info('reshifter.info/default-etcd: '+ endpoint);
    ibstorage.setItem('reshifter.info/default-apiversion', apiversion);
    console.info('reshifter.info/default-apiversion: '+ apiversion);
    ibstorage.setItem('reshifter.info/default-remote', remote);
    console.info('reshifter.info/default-remote: '+ remote);
    $('#config-result').html('<h2>Result</h2><div>All settings stored locally:</div><div><ul>');
    $('#config-result').append('<li>Using etcd endpoint: <code>' + ibstorage.getItem('reshifter.info/default-etcd') + '</code></li>');
    $('#config-result').append('<li>Using etcd API version: <code>' + ibstorage.getItem('reshifter.info/default-apiversion') + '</code></li>');
    $('#config-result').append('<li>Using backup target: <code>' + ibstorage.getItem('reshifter.info/default-remote') + '</code></li>');
    $('#config-result').append('</ul></div>');
    $('#config-result').append('<div style="margin-bottom: 50px"></div>');
  });

  $('#dobackup').click(function(event) {
    var ep = $('#endpoint').val();
    var filter = $('#filter').val();
    var default_remote = ibstorage.getItem('reshifter.info/default-remote');
    var default_apiversion = ibstorage.getItem('reshifter.info/default-apiversion');
    var sepidx = default_remote.indexOf(' ');
    var remote = '';
    var payload = '{"endpoint": "' + ep + '", "filter": "' + filter + '", "apiversion": "' + default_apiversion + '" }';
    if (sepidx != -1){
      remote =  default_remote.substring(sepidx+1, default_remote.lastIndexOf(' '));
      bucket =  default_remote.substring(default_remote.lastIndexOf(' ')+1);
      payload = '{"endpoint": "' + ep + '", "filter": "' + filter + '", "remote": "' + remote +'", "bucket": "' + bucket +'" }',
      console.info('Backing up to remote [' + remote + '] in bucket [' + bucket + ']');
    }
    console.info('Using filter [' + filter + ']');
    $('#backup-result').html('<div><img src="./img/standby.gif" alt="please wait" width="64px"></div>');
    $.ajax({
        type: "POST",
        url: 'http://localhost:8080/v1/backup',
        dataType: 'json',
        async: false,
        data: payload,
        error: function (d) {
          console.info(d);
          $('#backup-result').html('<h2>Result</h2>')
          $('#backup-result').append('<div>There was a problem carrying out the backup:<div><code>' + d.responseText + '</code></div></div>')
          $('#backup-result').append('<div style="margin-bottom: 50px"></div>');
        },
        success: function (d) {
          console.info(d);
          $('#backup-result').html('<h2>Result</h2>')
          if(d.outcome == 'success'){
            ibstorage.setItem('reshifter.info/last-backup-id', d.backupid)
            $('#backup-result').append('<div>The backup with ID <code>' + d.backupid +'</code> is now available <a href="/v1/backup/'+ d.backupid + '">here</a> for download.</div>')
          } else{
            $('#backup-result').append('<div>There was a problem carrying out the backup:<div><code>'+ d + '</code></div></div>')
          }
          $('#backup-result').append('<div style="margin-bottom: 50px"></div>');
        }
    })
  });

  $('#dolistbackups').click(function(event) {
    var default_remote = ibstorage.getItem('reshifter.info/default-remote');
    var sepidx = default_remote.indexOf(' ');
    var remote = '';
    var payload = '';
    if (sepidx != -1){ // we have remote backup selected
      remote =  default_remote.substring(sepidx+1, default_remote.lastIndexOf(' '));
      bucket =  default_remote.substring(default_remote.lastIndexOf(' ')+1);
      console.info('Listing backups from remote [' + remote + '] in bucket [' + bucket + ']');
      payload = '?remote=' + encodeURIComponent(remote) +'&bucket=' + encodeURIComponent(bucket);
    }
    $('#restore-result').html('<div><img src="./img/standby.gif" alt="please wait" width="64px"></div>');
    $.ajax({
        type: "GET",
        url: 'http://localhost:8080/v1/backup/all'+payload,
        dataType: 'json',
        error: function (d) {
          console.info(d);
          $('#restore-result').html('<h2>Result</h2>');
          $('#restore-result').append('<div>There was a problem listing available backups:<div><code>' + d.responseText + '</code></div></div>');
          $('#restore-result').append('<div style="margin-bottom: 50px"></div>');
        },
        success: function (d) {
          var backups = d.backupids;
          console.info(d);
          if (sepidx != -1){ // we have remote backup selected
            $('#restore-result').html('<h2>Result</h2><div>The following backups are available in the remote storage:</div>');
          } else {
            $('#restore-result').html('<h2>Result</h2><div>The following backups are available, locally:</div>');
          }
          for (var i = 0; i < backups.length; i++) {
            $('#restore-result').append('<div style="margin: 10px;"><code>' + backups[i] +'</code></div>');
          }
          $('#restore-result').append('<div style="margin-bottom: 50px"></div>');
        }
    })
  });

  $('#dorestore').click(function(event) {
    var ep = $('#endpoint').val();
    var bid = $('#backupid').val();
    var default_remote = ibstorage.getItem('reshifter.info/default-remote');
    var sepidx = default_remote.indexOf(' ')
    var remote = ''
    var payload = '{ "endpoint": "' + ep + '", "backupid": "' + bid +'" }'
    if (sepidx != -1){
      remote =  default_remote.substring(sepidx+1, default_remote.lastIndexOf(' '))
      bucket =  default_remote.substring(default_remote.lastIndexOf(' ')+1)
      payload = '{"endpoint": "' + ep + '", "backupid": "' + bid + '", "remote": "' + remote +'", "bucket": "' + bucket +'" }'
    }
    $('#restore-result').html('<div><img src="./img/standby.gif" alt="please wait" width="64px"></div>');
    $.ajax({
        type: "POST",
        url: 'http://localhost:8080/v1/restore',
        dataType: 'json',
        async: false,
        data: payload,
        error: function (d) {
          console.info(d);
          $('#restore-result').html('<h2>Result</h2>');
          $('#restore-result').append('<div>There was a problem carrying out the restore:<div><code>' + d.responseText + '</code></div></div>');
          $('#restore-result').append('<div style="margin-bottom: 50px"></div>');
        },
        success: function (d) {
          console.info(d);
          $('#restore-result').html('<h2>Result</h2>');
          if(d.outcome == 'success'){
            $('#restore-result').append('<div>Restored ' + d.keysrestored + ' keys from backup with ID <code>' + bid +'</code> to <code>'+ ep + '</code> in ' + d.elapsedtimeinsec + ' seconds.</div>');
          } else{
            $('#restore-result').append('<div>There was a problem carrying out the restore:<div><code>' + d + '</code></div></div>');
          }
          $('#restore-result').append('<div style="margin-bottom: 50px"></div>');
        }
    })
  });

  $('#doupload').click(function(event) {
    var bfdata = new FormData($('#bfuploader')[0]);
    $('#restore-result').html('<div><img src="./img/standby.gif" alt="please wait" width="64px"></div>');
    $.ajax({
        type: "POST",
        url: 'http://localhost:8080/v1/restore/upload',
        data: bfdata,
        cache: false,
        contentType: false,
        processData: false,
        error: function (d) {
          console.info(d);
          $('#restore-result').html('<h2>Result</h2>');
          $('#restore-result').append('<div>There was a problem uploading the backup file:<div><code>' + d.responseText + '</code></div></div>');
          $('#restore-result').append('<div style="margin-bottom: 50px"></div>');
        },
        success: function (d) {
          console.info(d);
          d = $.parseJSON(d);
          r = d.received;
          $('#restore-result').html('<h2>Result</h2>');
          $('#restore-result').append('<div>Received ' +  r + ' bytes in the backup file, now available via local storage for backup.');
          $('#restore-result').append('<div style="margin-bottom: 50px"></div>');
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

  $('.sidebar-close').click(function(){
    $('body').toggleClass('sidebar-opened').parent().toggleClass('sidebar-opened');
  });

  /* ---------- Disable moving to top ---------- */
  $('a[href="#"][data-top!=true]').click(function(e){
    e.preventDefault();
  });
});


// utils

function getVersion(){
  $.ajax({
      type: "GET",
      url: 'http://localhost:8080/v1/version',
      async: false,
      error: function (d) {
        console.info(d);
      },
      success: function (d) {
        console.info(d);
        $('.app-footer').html(d)
      }
  })
}

function setDefaults(){
  var default_ep = ibstorage.getItem('reshifter.info/default-etcd');
  var default_apiversion = ibstorage.getItem('reshifter.info/default-apiversion');
  var default_remote = ibstorage.getItem('reshifter.info/default-remote');
  var sepidx = 0;
  var remote = '';
  var bucket = '';

  if (default_ep == null) {
    ibstorage.setItem('reshifter.info/default-etcd', 'http://localhost:2379');
  }
  if (default_apiversion == null) {
    ibstorage.setItem('reshifter.info/default-apiversion', 'auto');
  }
  if (default_remote == null) {
    ibstorage.setItem('reshifter.info/default-remote', 's3 play.minio.io:9000 reshifter-xxx');
  }
  sepidx = default_remote.indexOf(' ');
  console.info('reshifter.info/default-etcd: '+default_ep)
  $('#endpoint').val(default_ep);

  console.info('reshifter.info/default-apiversion: '+default_apiversion)
  $('#apiversion').val(default_apiversion);

  console.info('reshifter.info/default-remote: '+default_remote);

  remote =  default_remote.substring(sepidx+1, default_remote.lastIndexOf(' '));
  if (sepidx == -1){ // download as ZIP file
    $("#remote[download=s3]").prop('checked', true);
  } else { // we have a remote configured
    $("#remote[value=s3]").prop('checked', true);
    bucket = default_remote.substring(default_remote.lastIndexOf(' ')+1);
    $('#bucket').val(bucket);
  }
}

function initBackup(){
  var default_ep = ibstorage.getItem('reshifter.info/default-etcd');
  var default_remote = ibstorage.getItem('reshifter.info/default-remote');
  var sepidx = default_remote.indexOf(' ');
  var remote = '';
  var bucket ='';

  $('#endpoint').val(default_ep);

  if (sepidx == -1){ // download as ZIP file
    console.info('User will download ZIP file');
  } else { // we have a remote configured
    remote =  default_remote.substring(sepidx+1, default_remote.lastIndexOf(' '));
    bucket =  default_remote.substring(default_remote.lastIndexOf(' ')+1);
    $('#backup-result').html('<div>Backing up to: <code>'+ remote +'</code>, into bucket <code>' + bucket + '</code></div>');
    console.info('Backing up to [' + remote + '] into bucket [' + bucket + ']');
  }
}

function initRestore(){
  var default_ep = ibstorage.getItem('reshifter.info/default-etcd');
  var default_remote = ibstorage.getItem('reshifter.info/default-remote');
  var sepidx = default_remote.indexOf(' ');
  var remote = '';
  var bucket ='';

  $('#endpoint').val(default_ep);

  if (sepidx == -1){ // upload from local storage
    console.info('User will use ZIP file from local storage');
  } else { // we have a remote configured
    remote =  default_remote.substring(sepidx+1, default_remote.lastIndexOf(' '));
    bucket =  default_remote.substring(default_remote.lastIndexOf(' ')+1);
    $('#restore-result').html('<div>Restoring from remote <code>'+ remote +'</code>, from bucket <code>' + bucket + '</code></div>');
    console.info('Restoring from remote [' + remote + '] from bucket [' + bucket + ']');
  }

  last_backup_id = ibstorage.getItem('reshifter.info/last-backup-id');
  if (last_backup_id !== '') {
    console.info('Using last backup ID: '+last_backup_id);
    $('#backupid').val(last_backup_id);
  }
}

function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

function init(url) {
  $('[rel="tooltip"],[data-rel="tooltip"]').tooltip({"placement":"bottom",delay: { show: 400, hide: 200 }});
  $('[rel="popover"],[data-rel="popover"],[data-toggle="popover"]').popover();
}
