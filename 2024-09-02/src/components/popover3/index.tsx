import React from 'react';
import * as Popover from '@radix-ui/react-popover';
import { Button, Flex, Avatar, Box, TextArea, Text, Checkbox } from '@radix-ui/themes';
import { ChatBubbleIcon } from '@radix-ui/react-icons';
import { keyframes } from '@emotion/react';
import { styled } from '@mui/system';

const slideUpAndFade = keyframes`
  from {
    opacity: 0;
    transform: translateY(2px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
`;

const slideDownAndFade = keyframes`
  from {
    opacity: 1;
    transform: translateY(0);
  }
  to {
    opacity: 0;
    transform: translateY(2px);
  }
`;

const StyledContent = styled(Popover.Content)`
  animation: ${slideUpAndFade} 800ms cubic-bezier(0.16, 1, 0.3, 1);
  &[data-state="closed"] {
    animation: ${slideDownAndFade} 600ms cubic-bezier(0.16, 1, 0.3, 1);
  }
`;

const PopoverDemo = () => {
  return (
    <Popover.Root>
      <Popover.Trigger>
        <Button variant="soft">
          <ChatBubbleIcon width="16" height="16" />
          Comment
        </Button>
      </Popover.Trigger>
      <StyledContent>
        <Flex gap="3">
          <Avatar
            size="2"
            src="https://images.unsplash.com/photo-1607346256330-dee7af15f7c5?&w=64&h=64&dpr=2&q=70&crop=focalpoint&fp-x=0.67&fp-y=0.5&fp-z=1.4&fit=crop"
            fallback="A"
            radius="full"
          />
          <Box flexGrow="1">
            <TextArea placeholder="Write a comment…" style={{ height: 80 }} />
            <Flex gap="3" mt="3" justify="between">
              <Flex align="center" gap="2" asChild>
                <Text as="label" size="2">
                  <Checkbox />
                  <Text>Send to group</Text>
                </Text>
              </Flex>

              <Popover.Close>
                <Button size="1">Comment</Button>
              </Popover.Close>
            </Flex>
          </Box>
        </Flex>
      </StyledContent>
    </Popover.Root>
  );
};

export default PopoverDemo;
