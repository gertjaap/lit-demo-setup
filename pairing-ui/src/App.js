import React, { Component } from 'react';
import './App.css';
import { ScaleLoader } from 'react-spinners';
import { Navbar, NavbarBrand } from 'reactstrap';
import {Row,Col} from 'reactstrap';

import QRCode from 'qrcode.react';
import { Table } from 'reactstrap';
import { Button } from 'reactstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faHandshake, faQrcode } from '@fortawesome/free-solid-svg-icons'
import { faLock } from '@fortawesome/free-solid-svg-icons'
import { faFileAlt } from '@fortawesome/free-solid-svg-icons'
import { faSync } from '@fortawesome/free-solid-svg-icons'
import { faTrashAlt } from '@fortawesome/free-solid-svg-icons'

class App extends Component {

  constructor(props) {
    super(props);
    
    this.state = {
      IsDropping: "",
      BlockHeights: {},
      PairingNode: null,
      PairedNodes: [],
      Nodes: [],
      step: "install",
      currentNode: "",
      currentNodeAddress: "",
    };

    this.update = this.update.bind(this);
    this.confirmInstall = this.confirmInstall.bind(this);
    this.confirmPair = this.confirmPair.bind(this);
    this.cleanupNodes = this.cleanupNodes.bind(this);
  }

  confirmInstall() {
    this.setState({step:"pair"});
  }

  confirmPair(pairedNodeName) {
    var paired = this.state.PairedNodes;
    paired.push({name:pairedNodeName, pairedDate: new Date()});
    this.approvePendingRequests(pairedNodeName);
    fetch("/api/nodes/new").then(() => { this.update(); });
    this.setState({step:"install", PairedNodes: paired});
  }

  approvePendingRequests(nodeName) {
    fetch("/api/nodes/pendingauth/" + nodeName)
    .then(res => res.json())
    .then(res => {
        res.forEach((key) => {
            fetch("/api/nodes/auth/" + nodeName + "/" + key + "/1")
            .then(res => res.json())
        });
    });
  }


  cleanupNodes() {
     for(var paired of this.state.PairedNodes) {
       if(paired.pairedDate.valueOf() < (new Date().valueOf()-600000) && paired.deleted !== true) {
         fetch("/api/nodes/delete/" + paired.name).then(() => { paired.deleted = true; }).catch((err) => { console.error(err); });
       }
     } 
  }

  update() { 
      fetch("/api/nodes/list")
      .then(res => res.json())
      .then(nodes => {

        fetch("/api/chain/height")
        .then(res => res.json())
        .then(heights => {

            this.setState({Nodes: nodes, BlockHeights: heights});
            this.cleanupNodes();
          
        }).catch(err => {
          
        });
        
      }).catch(err => {

      });
    
  }

  render() {

    var nextNode = "";
    var nextNodeName = "";
    var blockHeightNames = {};
    blockHeightNames[257] = "bitcoind";
    blockHeightNames[258] = "litecoind";
    blockHeightNames[262] = "dummyusdd";

    if(!(this.state.Nodes === null || this.state.Nodes.length == 0)) {
      for(var node of this.state.Nodes) {
        var used = false;
        for(var pairedNode of this.state.PairedNodes) {
          if(pairedNode.name === node.Name) { 
            used = true;
            break;
          }
        }

        if(used) {
          console.log("Node is already in use: ", node.Name);
          continue;
        }
        // make sure blockheights are okay
        var blockHeightsOkay = true;
        for(var balance of node.Balances) {
          var curHeight = this.state.BlockHeights[blockHeightNames[balance.CoinType]];

          if(!(balance.SyncHeight >= curHeight-1)) {
            console.log("Mismatching blockheight for ", balance.CoinType, ": ", balance.SyncHeight, curHeight);
            blockHeightsOkay = false;
            break;
          }
        }

        if(!blockHeightsOkay) {
          console.log("Node has non matching blockheights:", node.Name);
          continue;
        }

        var balancesOkay = true; 
        for(var balance of node.Balances) {
          if(balance.CoinType === 257 && balance.ChanTotal < 400000) balancesOkay = false;
          if(balance.CoinType === 258 && balance.ChanTotal < 50000000) balancesOkay = false;
          if(balance.CoinType === 262 && balance.ChanTotal < 1000000000) balancesOkay = false;
          
        }

        if(!balancesOkay) {
          console.log("Node has non matching balances:", node.Name);
          continue;
        }

        nextNode = node.Address + "@dcidemo.media.mit.edu:" + node.PublicLitPort;
        nextNodeName = node.Name;
      }
    }
    
    var mainContent = (<div>
      <h1>To get started, install the lit mobile app:</h1>
          
          <Row>
            <Col xs={1}>&nbsp;</Col>
            <Col xs={3}>
              <h2>Android:</h2>
              <QRCode size={192} value="https://play.google.com/apps/testing/edu.mit.dci.lit " />
            </Col>
            <Col xs={4}>&nbsp;</Col>
            <Col xs={3}>
              <h2>iOS:</h2>
              <QRCode size={192} value="https://testflight.apple.com/join/AmvSmBXO" />
            </Col>
            <Col xs={1}>&nbsp;</Col>
          </Row>
  
          <br/> &nbsp; <br/>
          <center>
            <Button onClick={this.confirmInstall}>
                Next &gt;
            </Button>
          </center>
    </div>);

    var pair = <div>No node available</div>

    if (nextNode !== "") {
      pair = (<center>
          <QRCode size={192} value={nextNode} />
        </center>)
    }

    if(this.state.step === "pair") {
      mainContent = (<div>

        <h1>Use the QR code below to pair:</h1>
        
        {pair}

        <br/> &nbsp; <br/>
          <center>
            <Button onClick={() => { this.confirmPair(nextNodeName) }}>
                Next &gt;
            </Button>
          </center>
      </div>)
    }

    return (
      <div className="App">
        
        {mainContent}

      </div>
    );
  }

  componentDidMount() {
    this.update();
    setInterval(() => {this.update()}, 10000);
  }
}

export default App;

/**
 * 
 * 
 * <CssBaseline />
        <AppBar position="static">
          <Toolbar>
            <Typography variant="title" color="inherit">
              LIT Demo Environment - {blockHeights}
            </Typography>
            <IconButton onClick={() => this.mineBlock()} aria-label="Mine">
                      <GavelIcon />
            </IconButton>
          </Toolbar>
        </AppBar>



        <Grid container spacing={24}>
          {this.state.Nodes.map(n => {
            return (
              <Grid key={n.Name} item xs={12} sm={6}>
                <Card>
                  <CardContent>
                    <Typography variant="headline" component="h1" color="textSecondary">
                      {n.Name}
                    </Typography>
                    <Typography component="small">
                    <b>Public node address:</b><IconButton onClick={() => this.copyPubNodeToClipboard(n)} aria-label="Copy">
                      <ContentCopy />
                    </IconButton><br/>
                      {n.Address + '@' + document.location.hostname + ':' + n.PublicLitPort}<br/>
                    </Typography>
                    <Typography component="small">
                    <b>Private node address:</b><IconButton onClick={() => this.copyPrivNodeToClipboard(n)} aria-label="Copy">
                      <ContentCopy />
                    </IconButton><br/>
                    {n.Address + '@' + n.Name}
                    </Typography>
                    <Typography component="small">
                    <b>Balance:</b><br/>
                    {(n.Balances[257] / 100000) + ' mBTC-R'}
                    </Typography>
                  </CardContent>
                  <CardActions>
                    <IconButton aria-label="Delete" onClick={() => this.dropNode(n)}>
                      <DeleteIcon  />
                    </IconButton>
                    <IconButton aria-label="Fund" onClick={() => this.fundNode(n)}>
                      <AttachMoneyIcon />
                    </IconButton>
                    <IconButton aria-label="Open LIT UI"  onClick={() => this.openLitUIForNode(n)}>
                      <OpenInNewIcon />
                    </IconButton>
                  </CardActions>
                </Card>
              </Grid>
            );
          })}
          
          <Grid key="New" item xs={12} sm={6}>
            {creation}
            <Button onClick={this.newNode.bind(this)}  disabled={this.state.IsCreating} variant="fab" color="primary" aria-label="add">
              <AddIcon />
            </Button>
          </Grid>
        </Grid>
        
 */