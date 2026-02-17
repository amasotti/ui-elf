import React from 'react';
import { Button, Dialog } from '@mui/material';

function MyComponent() {
  return (
    <div>
      <Button>Click me</Button>
      <Dialog open={true}>
        <p>Hello</p>
      </Dialog>
    </div>
  );
}

export default MyComponent;
