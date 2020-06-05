import { Router, route } from 'preact-router';
import { authKey } from './const/storage';
import * as pages from './const/pages';
import { Component } from 'preact';

import Home from './routes/home';
import Signin from './routes/signin';

import './style';

export default class App extends Component {
    handleRoute = async e => {
        this.setState({
            currentUrl: e.url
        });
    }
    constructor() {
        super();
        this.state = { currentUrl: '' };
    }
    componentDidMount() {
        const userKey = localStorage.getItem(authKey);
        console.log(`userKey: ${userKey}`);
        if (!userKey) {
            route(pages.signin, true);
        }
    }
    render() {
        return (
            <div id="app">
                <Router onChange={this.handleRoute}>
                    <Home path={pages.home} />
                    <Signin path={pages.signin} />
                    <Home default />
                </Router>
            </div>
        );
    }
}
