package main

const exhibitTemplate = `
<html>
  <head>
    <title>{{.Header}}</title>
    

    <script src="http://trunk.simile-widgets.org/exhibit/api/exhibit-api.js"></script>
    <link href="allSoftware.js"  rel="exhibit/data" />

    <style>
      body {
      margin: 0;
      padding: 0;
      font-family: "Lucida Grande","Tahoma","Helvetica","Arial",sans-serif;
      color: #222;
      }

      table, tr, td {
      font-size: inherit;
      }

      tr, td {
      vertical-align: top;
      }

      img, a img {
      border: none;
      }

      #main-content { background: white; }

      #title-panel { padding: 0.25in 0.5in; }

      #top-panels {
      padding: 0.5em 0.5in;
      border-top: 1px solid #BCB79E;
      border-bottom: 1px solid #BCB79E;
      background: #FBF4D3;
      }

      .exhibit-tileView-body { list-style: none; }

      .exhibit-collectionView-group-count { display: none; }

      div.name {
      font-weight: bold;
      font-size: 120%;
      }

      table.software {
      border: 1px solid #ddd;
      padding: 1em;
      margin: 0.5em 0;
      display: block;
      valign: top;
      }

      span.name {
      font-weight: bold;
      font-size: 120%;
      }
      
    </style>

    
  </head> 

  <body>
    <h1>Software</h1>
    Made with <a href="https://github.com/AAFC-MBB/rpmout">rpmout</a> and <a href="http://www.simile-widgets.org/exhibit/">Exhibit</a>
    <table width="100%">
      <tr valign="top">
	<td ex:role="viewPanel">
	  <div ex:role="view" ex:label="List" >></div>
	  <!-- BEGIN LENS -->
	  <table ex:role="lens" class="software" style="display: none;"><>
            <tr>
              <td>
		<div>
		  <b><em><span class="name" ex:content=".label" ></span></em></b>:
		  <span ex:content=".ShortDescription" ></span>
		  <p>
		</div>
		<div><b>Description:</b>
		<span ex:content=".Description" ></span>
		</div>
		<div>
		<b>License:</b>
		<span ex:content=".License" ></span>
		</div>
		<div>
		<b>Group:</b>
		<span ex:content=".Group" ></span>
		</div>
		<div>
		  <b>URL:</b>
		  <a ex:href-content=".Url"><span ex:content=".Url" ></span></a>
		</div>
              </td>
            </tr>
	  </table>
	  <!-- END LENS -->
	</td>

	<td width="35%">
	  <b>Search</b>
	  <div ex:role="facet" ex:facetClass="TextSearch"></div>
	  <div ex:role="facet" ex:expression=".Group" ex:facetLabel="Group" ex:height="19em">Software Group</div>
	  <div ex:role="facet" ex:sortMode="count" ex:expression=".License" ex:facetLabel="License"  ex:height="19em">License</div>
	  <div ex:role="facet" ex:expression=".Type" ex:facetLabel="Type" ex:height="4em">Software Type</div>
        </td>
      </tr>
    </table>
  </body>
  
</html>
`
