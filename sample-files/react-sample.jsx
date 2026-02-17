import React from 'react';
import { Button, Dialog, Form } from '@mui/material';

function ReactSample() {
  const [open, setOpen] = React.useState(false);

  return (
    <div>
      <Form onSubmit={(e) => e.preventDefault()}>
        <Button type="submit">Submit</Button>
      </Form>
      <Dialog open={open} onClose={() => setOpen(false)}>
        <div>Dialog content</div>
      </Dialog>
    </div>
  );
}

export default ReactSample;
