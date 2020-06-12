import { Component } from 'preact';
import {
    Input,
    InputAdornment,
    IconButton,
    InputLabel,
    TextField,
    Button,
    Container,
    Typography,
    Avatar,
} from '@material-ui/core';

import { FormControl, FormControlLabel, Checkbox } from '@material-ui/core'

import { Visibility, VisibilityOff, LockOutlined } from '@material-ui/icons';


export default class Signin extends Component {
    constructor() {
        super();
        this.state = {
            showPassword: false,
            username: null,
            password: null
        };
    }
    handleClickShowPassword = () => {
        let showPassword = this.state.showPassword;
        this.setState({ showPassword: !showPassword });
    }
    render() {
        return (
            <Container component="main" maxWidth="xs">
                <Avatar>
                    <LockOutlined />
                </Avatar>
                <Typography component="h1" variant="h5">
                    Sign in
                </Typography>
                <TextField
                    margin="normal"
                    required
                    fullWidth
                    label="Email Address"
                    name="email"
                    autoComplete="email"
                    autoFocus
                />
                <TextField
                    margin="normal"
                    required
                    fullWidth
                    name="password"
                    label="Password"
                    type={this.state.showPassword ? "text" : "password"}
                    autoComplete="current-password"
                    InputProps={{
                        endAdornment: (
                            <InputAdornment onClick={this.handleClickShowPassword}>
                                {this.state.showPassword ? <Visibility /> : <VisibilityOff />}
                            </InputAdornment>
                        ),
                    }}
                />
                <Button
                    type="submit"
                    fullWidth
                    variant="contained"
                    color="primary"
                >
                    Sign In
                    </Button>
            </Container>
        );
    }
}
