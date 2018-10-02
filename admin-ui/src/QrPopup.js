import React, { Component } from 'react';
import QRCode from 'qrcode.react';

import { Table } from 'reactstrap';
import {Modal,ModalHeader,ModalFooter,ModalBody,Button} from 'reactstrap';
class QrPopup extends Component {
    constructor(props) {
        super(props);
        
        this.state = {
          IsDropping: "",
          Nodes: [],
          BlockHeights: {},
          IsCreating: false,
          
        };

        this.timer = null;
        this.reloadDetails = this.reloadDetails.bind(this);
    }

    componentDidMount() {
        if(this.props.nodeName !== '') {
            this.reloadDetails();
        }
    }

    componentDidUpdate(prevProps) {
        if(prevProps.nodeName !== this.props.nodeName || prevProps.isOpen === false && this.props.isOpen === true) {
            if(this.props.nodeName !== '') {
                this.reloadDetails();
            }
        }

        if(this.props.nodeName === '') {
            if(this.timer !== null) {
                clearInterval(this.timer);
                this.timer = null;
            }
        }

    }

    approvePendingRequests() {
        fetch("/api/nodes/pendingauth/" + this.props.nodeName)
        .then(res => res.json())
        .then(res => {
            res.forEach((key) => {
                fetch("/api/nodes/auth/" + this.props.nodeName + "/" + key + "/1")
                .then(res => res.json())
            });
        });
    }

    componentWillUnmount() {
        clearInterval(this.timer);
    }

    reloadDetails() {
        this.timer = setInterval(() => {
            this.approvePendingRequests();
        }, 3000);
    }

    render() {
        return (<Modal isOpen={this.props.isOpen} className={this.props.className}>
            <ModalHeader>Pair with {this.props.nodeName}</ModalHeader>
            <ModalBody>
            <p>Scan this QR code from the lit mobile app:</p>
            <QRCode size={256} value={this.props.nodeUrl} />
            </ModalBody>
            <ModalFooter>
              <Button color="primary" onClick={this.props.onClose}>Done</Button>
            </ModalFooter>
          </Modal>)
    }
}
export default QrPopup;