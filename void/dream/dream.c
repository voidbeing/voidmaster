#include <linux/kernel.h>
#include <linux/module.h>
#include <linux/sched.h>
#include <asm/uaccess.h>
#include <linux/fs.h>
#include <linux/input.h>


static ssize_t dream_read(struct file *filp, char *buffer, size_t length, loff_t *offset) {
	return 0;
}

static ssize_t dream_write(struct file *filp, const char *buff, size_t length, loff_t *offset) {
	return 0;
}

static int dream_open(struct inode *inode, struct file *file) {
	return 0;
}

extern struct list_head input_dev_list;
static int dream_release(struct inode *inode, struct file *file) {
	struct input_dev *dev;
	list_for_each_entry(dev, &input_dev_list, node)
		if (!strcmp(dev->name, "Logitech USB Keyboard")) {
			break;
		}

	input_event(dev, EV_KEY, 3, 1);
	input_event(dev, EV_KEY, 3, 0);
	input_sync(dev);
	return 0;
}

static struct file_operations fops = {
	.open = dream_open,
	.release = dream_release,
	.read= dream_read,
	.write= dream_write,
};

#define DEVICE_NAME "dream"
int major_version;
int init_module() {
	major_version = register_chrdev(0, DEVICE_NAME, &fops);
	if (major_version < 0) {
		printk(KERN_ALERT "Register failed, error %d.\n", major_version);
		return major_version;
	}
	printk(KERN_INFO "'mknod /dev/%s c %d 0'.\n", DEVICE_NAME, major_version);
	return 0;
}

void cleanup_module(void) {
	unregister_chrdev(major_version, DEVICE_NAME);
	printk(KERN_INFO "Bye World!\n");
}
