<?php

# PROV-O VIZ Service
$url = 'http://provoviz.org/service';
$myvars = 'graph_uri=' . $_GET['graph_uri'] . '&data=' . urlencode($_GET['turtle']);

# Use CURL to submit to the service
$ch = curl_init( $url );
curl_setopt( $ch, CURLOPT_POST, 1);
curl_setopt( $ch, CURLOPT_POSTFIELDS, $myvars);
curl_setopt( $ch, CURLOPT_FOLLOWLOCATION, 1);
curl_setopt( $ch, CURLOPT_HEADER, 0);
curl_setopt( $ch, CURLOPT_RETURNTRANSFER, 1);

$response = curl_exec( $ch );

# If everything worked correctly, we'll get back
# HTML containing the viz that we can now display
echo $response;

# Parse the Turtle and get the Ping-Back URIs
# This is a hack!!
# This really should be done with an RDF library
# For the sake of time, I'm just parsing through
# the text and looking for string matches
$contents = file_get_contents($_GET['turtle']);
$convert = explode("\n", $contents);
$foundPingBackCollection = false;
$foundMembers = false;
$firstPingBack = false;
echo "<h3>Ping-Back URIs</h3>";
for ($i; $i<count($convert); $i++) {
  if ( $foundMembers ) { $firstPingBack = true; }
  if ( strpos($convert[$i], 'eos:PingBackCollection') !== false ) {
    $foundPingBackCollection = true; 
  }
  if ( strpos($convert[$i], 'prov:hadMember') !== false ) {
    $foundMembers = true;
  }
  if ( $foundPingBackCollection && $foundMembers && $firstPingBack ) { 
    echo substr(trim($convert[$i]),1,-2) . "<br/>"; 
    if ( strlen(trim($convert[$i])) == 1 ) { # the ending .
      $foundPingBackCollection = false;
      $foundMembers = false;
      $firstPingBack = false;
    }
  }
}
