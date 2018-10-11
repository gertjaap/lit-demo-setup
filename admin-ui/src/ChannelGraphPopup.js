import React, { Component } from 'react';
import QRCode from 'qrcode.react';

import { Table } from 'reactstrap';
import {Modal,ModalHeader,ModalFooter,ModalBody,Button} from 'reactstrap';
class ChannelGraphPopup extends Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (<Modal isOpen={this.props.isOpen} className={this.props.className} style={{maxWidth:'100%', width:'90%'}}>
            <ModalHeader>Channel Graph</ModalHeader>
            <ModalBody>
                <center>
            <img style={{maxWidth:'100%'}} src="/api/nodes/graph" />
            </center>
            </ModalBody>
            <ModalFooter>
              <Button color="primary" onClick={this.props.onClose}>Done</Button>
            </ModalFooter>
          </Modal>)
    }
}
export default ChannelGraphPopup;